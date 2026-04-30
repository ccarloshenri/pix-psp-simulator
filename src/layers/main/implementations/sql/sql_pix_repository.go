package sqlrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

// PixRepository persists and retrieves Pix payments in Postgres.
type PixRepository struct {
	db *sql.DB
}

func NewPixRepository(db *sql.DB) *PixRepository {
	return &PixRepository{db: db}
}

func (r *PixRepository) Save(pix models.Pix) error {
	devolucoes := pix.Devolucoes
	if devolucoes == nil {
		devolucoes = []models.Devolucao{}
	}
	devolucoesJSON, err := json.Marshal(devolucoes)
	if err != nil {
		return fmt.Errorf("marshal devolucoes: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO pix_payments (end_to_end_id, txid, valor, horario, infopagador, devolucoes)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		pix.EndToEndID, pix.TxID, pix.Valor, pix.Horario,
		nullableString(pix.Infopagador), devolucoesJSON,
	)
	return err
}

func (r *PixRepository) FindByE2EID(e2eid string) (*models.Pix, error) {
	row := r.db.QueryRow(`
		SELECT end_to_end_id, txid, valor, horario, infopagador, devolucoes
		FROM pix_payments WHERE end_to_end_id = $1`, e2eid)
	return r.scanRow(row)
}

func (r *PixRepository) FindByTxID(txid string) ([]models.Pix, error) {
	rows, err := r.db.Query(`
		SELECT end_to_end_id, txid, valor, horario, infopagador, devolucoes
		FROM pix_payments WHERE txid = $1`, txid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PixRepository) FindAll(filters interfaces.PixFilters) ([]models.Pix, error) {
	query := `
		SELECT end_to_end_id, txid, valor, horario, infopagador, devolucoes
		FROM pix_payments WHERE 1=1`

	args := make([]any, 0)
	argIndex := 1

	if filters.TxID != "" {
		query += fmt.Sprintf(" AND txid = $%d", argIndex)
		args = append(args, filters.TxID)
		argIndex++
	}
	if filters.Inicio != "" {
		query += fmt.Sprintf(" AND horario >= $%d", argIndex)
		args = append(args, filters.Inicio)
		argIndex++
	}
	if filters.Fim != "" {
		query += fmt.Sprintf(" AND horario <= $%d", argIndex)
		args = append(args, filters.Fim)
		argIndex++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PixRepository) AddDevolucao(e2eid string, dev models.Devolucao) error {
	pix, err := r.FindByE2EID(e2eid)
	if err != nil {
		return err
	}
	if pix == nil {
		return fmt.Errorf("pagamento não encontrado")
	}

	pix.Devolucoes = append(pix.Devolucoes, dev)
	return r.persistDevolucoes(e2eid, pix.Devolucoes)
}

func (r *PixRepository) UpdateDevolucao(e2eid string, dev models.Devolucao) error {
	pix, err := r.FindByE2EID(e2eid)
	if err != nil {
		return err
	}
	if pix == nil {
		return fmt.Errorf("pagamento não encontrado")
	}

	for i, d := range pix.Devolucoes {
		if d.ID == dev.ID {
			pix.Devolucoes[i] = dev
			return r.persistDevolucoes(e2eid, pix.Devolucoes)
		}
	}
	return fmt.Errorf("devolução não encontrada")
}

func (r *PixRepository) FindDevolucao(e2eid, devID string) (*models.Devolucao, error) {
	pix, err := r.FindByE2EID(e2eid)
	if err != nil {
		return nil, err
	}
	if pix == nil {
		return nil, nil
	}
	for _, d := range pix.Devolucoes {
		if d.ID == devID {
			return &d, nil
		}
	}
	return nil, nil
}

func (r *PixRepository) persistDevolucoes(e2eid string, devolucoes []models.Devolucao) error {
	devolucoesJSON, err := json.Marshal(devolucoes)
	if err != nil {
		return fmt.Errorf("marshal devolucoes: %w", err)
	}
	_, err = r.db.Exec(
		`UPDATE pix_payments SET devolucoes = $1 WHERE end_to_end_id = $2`,
		devolucoesJSON, e2eid,
	)
	return err
}

func (r *PixRepository) scanRow(row *sql.Row) (*models.Pix, error) {
	var (
		e2eid, txid, valor string
		horario            time.Time
		infopagador        sql.NullString
		devolucoesJSON     []byte
	)

	err := row.Scan(&e2eid, &txid, &valor, &horario, &infopagador, &devolucoesJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return buildPix(e2eid, txid, valor, horario, infopagador, devolucoesJSON)
}

func (r *PixRepository) scanRows(rows *sql.Rows) ([]models.Pix, error) {
	result := make([]models.Pix, 0)
	for rows.Next() {
		var (
			e2eid, txid, valor string
			horario            time.Time
			infopagador        sql.NullString
			devolucoesJSON     []byte
		)

		if err := rows.Scan(&e2eid, &txid, &valor, &horario, &infopagador, &devolucoesJSON); err != nil {
			return nil, err
		}

		pix, err := buildPix(e2eid, txid, valor, horario, infopagador, devolucoesJSON)
		if err != nil {
			return nil, err
		}
		result = append(result, *pix)
	}
	return result, rows.Err()
}

func buildPix(
	e2eid, txid, valor string,
	horario time.Time,
	infopagador sql.NullString,
	devolucoesJSON []byte,
) (*models.Pix, error) {
	pix := models.Pix{
		EndToEndID:  e2eid,
		TxID:        txid,
		Valor:       valor,
		Horario:     horario,
		Infopagador: infopagador.String,
		Devolucoes:  []models.Devolucao{},
	}

	if len(devolucoesJSON) > 0 {
		if err := json.Unmarshal(devolucoesJSON, &pix.Devolucoes); err != nil {
			return nil, fmt.Errorf("unmarshal devolucoes: %w", err)
		}
	}

	return &pix, nil
}
