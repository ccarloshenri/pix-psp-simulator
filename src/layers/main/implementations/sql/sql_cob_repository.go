package sqlrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

// CobRepository persists and retrieves Cob records in Postgres.
type CobRepository struct {
	db *sql.DB
}

func NewCobRepository(db *sql.DB) *CobRepository {
	return &CobRepository{db: db}
}

func (r *CobRepository) Save(cob models.Cob) error {
	devedorJSON, err := marshalNullable(cob.Devedor)
	if err != nil {
		return fmt.Errorf("marshal devedor: %w", err)
	}
	locJSON, err := marshalNullable(cob.Loc)
	if err != nil {
		return fmt.Errorf("marshal loc: %w", err)
	}
	infoJSON, err := json.Marshal(cob.InfoAdicionais)
	if err != nil {
		return fmt.Errorf("marshal info_adicionais: %w", err)
	}
	retiradaJSON, err := marshalNullable(cob.Valor.Retirada)
	if err != nil {
		return fmt.Errorf("marshal retirada: %w", err)
	}
	pixListJSON, err := json.Marshal(cob.Pix)
	if err != nil {
		return fmt.Errorf("marshal pix_list: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO cobs (
			txid, status, chave, revisao, expiracao, criacao,
			valor_original, modalidade_alteracao, location, pix_copia_e_cola,
			solicitacao_pagador, devedor, loc, info_adicionais, retirada, pix_list
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		cob.TxID, cob.Status, cob.Chave, cob.Revisao,
		cob.Calendario.Expiracao, cob.Calendario.Criacao,
		cob.Valor.Original, cob.Valor.ModalidadeAlteracao,
		nullableString(cob.Location), nullableString(cob.PixCopiaECola),
		nullableString(cob.SolicitacaoPagador),
		devedorJSON, locJSON, infoJSON, retiradaJSON, pixListJSON,
	)
	return err
}

func (r *CobRepository) FindByTxID(txid string) (*models.Cob, error) {
	row := r.db.QueryRow(`
		SELECT txid, status, chave, revisao, expiracao, criacao,
		       valor_original, modalidade_alteracao, location, pix_copia_e_cola,
		       solicitacao_pagador, devedor, loc, info_adicionais, retirada, pix_list
		FROM cobs WHERE txid = $1`, txid)
	return r.scanRow(row)
}

func (r *CobRepository) FindAll(filters interfaces.CobFilters) ([]models.Cob, error) {
	query := `
		SELECT txid, status, chave, revisao, expiracao, criacao,
		       valor_original, modalidade_alteracao, location, pix_copia_e_cola,
		       solicitacao_pagador, devedor, loc, info_adicionais, retirada, pix_list
		FROM cobs WHERE 1=1`

	args := make([]any, 0)
	argIndex := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}
	if filters.Inicio != "" {
		query += fmt.Sprintf(" AND criacao >= $%d", argIndex)
		args = append(args, filters.Inicio)
		argIndex++
	}
	if filters.Fim != "" {
		query += fmt.Sprintf(" AND criacao <= $%d", argIndex)
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

func (r *CobRepository) Update(cob models.Cob) error {
	devedorJSON, err := marshalNullable(cob.Devedor)
	if err != nil {
		return fmt.Errorf("marshal devedor: %w", err)
	}
	locJSON, err := marshalNullable(cob.Loc)
	if err != nil {
		return fmt.Errorf("marshal loc: %w", err)
	}
	infoJSON, err := json.Marshal(cob.InfoAdicionais)
	if err != nil {
		return fmt.Errorf("marshal info_adicionais: %w", err)
	}
	retiradaJSON, err := marshalNullable(cob.Valor.Retirada)
	if err != nil {
		return fmt.Errorf("marshal retirada: %w", err)
	}
	pixListJSON, err := json.Marshal(cob.Pix)
	if err != nil {
		return fmt.Errorf("marshal pix_list: %w", err)
	}

	result, err := r.db.Exec(`
		UPDATE cobs SET
			status = $1, chave = $2, revisao = $3, expiracao = $4, criacao = $5,
			valor_original = $6, modalidade_alteracao = $7, location = $8,
			pix_copia_e_cola = $9, solicitacao_pagador = $10,
			devedor = $11, loc = $12, info_adicionais = $13, retirada = $14, pix_list = $15
		WHERE txid = $16`,
		cob.Status, cob.Chave, cob.Revisao, cob.Calendario.Expiracao, cob.Calendario.Criacao,
		cob.Valor.Original, cob.Valor.ModalidadeAlteracao,
		nullableString(cob.Location), nullableString(cob.PixCopiaECola),
		nullableString(cob.SolicitacaoPagador),
		devedorJSON, locJSON, infoJSON, retiradaJSON, pixListJSON,
		cob.TxID,
	)
	if err != nil {
		return err
	}
	return requireOneRowAffected(result, "cobrança", cob.TxID)
}

func (r *CobRepository) Delete(txid string) error {
	result, err := r.db.Exec(`DELETE FROM cobs WHERE txid = $1`, txid)
	if err != nil {
		return err
	}
	return requireOneRowAffected(result, "cobrança", txid)
}

func (r *CobRepository) scanRow(row *sql.Row) (*models.Cob, error) {
	var (
		txid, status, chave, valorOriginal        string
		revisao, expiracao, modalidade            int
		location, pixCopiaECola, solPagador       sql.NullString
		criacao                                   time.Time
		devedorJSON, locJSON, infoJSON            []byte
		retiradaJSON, pixListJSON                 []byte
	)

	err := row.Scan(
		&txid, &status, &chave, &revisao, &expiracao, &criacao,
		&valorOriginal, &modalidade, &location, &pixCopiaECola, &solPagador,
		&devedorJSON, &locJSON, &infoJSON, &retiradaJSON, &pixListJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return buildCob(
		txid, status, chave, revisao, expiracao, criacao,
		valorOriginal, modalidade, location, pixCopiaECola, solPagador,
		devedorJSON, locJSON, infoJSON, retiradaJSON, pixListJSON,
	)
}

func (r *CobRepository) scanRows(rows *sql.Rows) ([]models.Cob, error) {
	result := make([]models.Cob, 0)
	for rows.Next() {
		var (
			txid, status, chave, valorOriginal  string
			revisao, expiracao, modalidade      int
			location, pixCopiaECola, solPagador sql.NullString
			criacao                             time.Time
			devedorJSON, locJSON, infoJSON      []byte
			retiradaJSON, pixListJSON           []byte
		)

		if err := rows.Scan(
			&txid, &status, &chave, &revisao, &expiracao, &criacao,
			&valorOriginal, &modalidade, &location, &pixCopiaECola, &solPagador,
			&devedorJSON, &locJSON, &infoJSON, &retiradaJSON, &pixListJSON,
		); err != nil {
			return nil, err
		}

		cob, err := buildCob(
			txid, status, chave, revisao, expiracao, criacao,
			valorOriginal, modalidade, location, pixCopiaECola, solPagador,
			devedorJSON, locJSON, infoJSON, retiradaJSON, pixListJSON,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, *cob)
	}
	return result, rows.Err()
}

func buildCob(
	txid, status, chave string,
	revisao, expiracao int,
	criacao time.Time,
	valorOriginal string,
	modalidade int,
	location, pixCopiaECola, solPagador sql.NullString,
	devedorJSON, locJSON, infoJSON, retiradaJSON, pixListJSON []byte,
) (*models.Cob, error) {
	cob := models.Cob{
		TxID:    txid,
		Status:  status,
		Chave:   chave,
		Revisao: revisao,
		Calendario: models.Calendario{
			Criacao:   criacao,
			Expiracao: expiracao,
		},
		Valor: models.CobValor{
			Original:            valorOriginal,
			ModalidadeAlteracao: modalidade,
		},
		Location:           location.String,
		PixCopiaECola:      pixCopiaECola.String,
		SolicitacaoPagador: solPagador.String,
	}

	if len(devedorJSON) > 0 {
		var devedor models.Devedor
		if err := json.Unmarshal(devedorJSON, &devedor); err != nil {
			return nil, fmt.Errorf("unmarshal devedor: %w", err)
		}
		cob.Devedor = &devedor
	}
	if len(locJSON) > 0 {
		var loc models.Loc
		if err := json.Unmarshal(locJSON, &loc); err != nil {
			return nil, fmt.Errorf("unmarshal loc: %w", err)
		}
		cob.Loc = &loc
	}
	if len(infoJSON) > 0 {
		if err := json.Unmarshal(infoJSON, &cob.InfoAdicionais); err != nil {
			return nil, fmt.Errorf("unmarshal info_adicionais: %w", err)
		}
	}
	if len(retiradaJSON) > 0 {
		var retirada models.Retirada
		if err := json.Unmarshal(retiradaJSON, &retirada); err != nil {
			return nil, fmt.Errorf("unmarshal retirada: %w", err)
		}
		cob.Valor.Retirada = &retirada
	}
	if len(pixListJSON) > 0 {
		if err := json.Unmarshal(pixListJSON, &cob.Pix); err != nil {
			return nil, fmt.Errorf("unmarshal pix_list: %w", err)
		}
	}

	return &cob, nil
}
