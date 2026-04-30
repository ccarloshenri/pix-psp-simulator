package sqlrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

// CobVRepository persists and retrieves CobV records in Postgres.
type CobVRepository struct {
	db *sql.DB
}

func NewCobVRepository(db *sql.DB) *CobVRepository {
	return &CobVRepository{db: db}
}

func (r *CobVRepository) Save(cobv models.CobV) error {
	devedorJSON, err := json.Marshal(cobv.Devedor)
	if err != nil {
		return fmt.Errorf("marshal devedor: %w", err)
	}
	locJSON, err := marshalNullable(cobv.Loc)
	if err != nil {
		return fmt.Errorf("marshal loc: %w", err)
	}
	infoJSON, err := json.Marshal(cobv.InfoAdicionais)
	if err != nil {
		return fmt.Errorf("marshal info_adicionais: %w", err)
	}
	multaJSON, err := marshalNullable(cobv.Valor.Multa)
	if err != nil {
		return fmt.Errorf("marshal valor_multa: %w", err)
	}
	jurosJSON, err := marshalNullable(cobv.Valor.Juros)
	if err != nil {
		return fmt.Errorf("marshal valor_juros: %w", err)
	}
	abatimentoJSON, err := marshalNullable(cobv.Valor.Abatimento)
	if err != nil {
		return fmt.Errorf("marshal valor_abatimento: %w", err)
	}
	descontoJSON, err := marshalNullable(cobv.Valor.Desconto)
	if err != nil {
		return fmt.Errorf("marshal valor_desconto: %w", err)
	}
	pixListJSON, err := json.Marshal(cobv.Pix)
	if err != nil {
		return fmt.Errorf("marshal pix_list: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO cobvs (
			txid, status, chave, revisao, data_de_vencimento, validade_apos_vencimento,
			criacao, valor_original, pix_copia_e_cola, solicitacao_pagador,
			devedor, loc, info_adicionais, valor_multa, valor_juros,
			valor_abatimento, valor_desconto, pix_list
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		cobv.TxID, cobv.Status, cobv.Chave, cobv.Revisao,
		cobv.Calendario.DataDeVencimento, cobv.Calendario.ValidadeAposVencimento,
		cobv.Calendario.Criacao,
		cobv.Valor.Original,
		nullableString(cobv.PixCopiaECola), nullableString(cobv.SolicitacaoPagador),
		devedorJSON, locJSON, infoJSON,
		multaJSON, jurosJSON, abatimentoJSON, descontoJSON,
		pixListJSON,
	)
	return err
}

func (r *CobVRepository) FindByTxID(txid string) (*models.CobV, error) {
	row := r.db.QueryRow(`
		SELECT txid, status, chave, revisao, data_de_vencimento, validade_apos_vencimento,
		       criacao, valor_original, pix_copia_e_cola, solicitacao_pagador,
		       devedor, loc, info_adicionais, valor_multa, valor_juros,
		       valor_abatimento, valor_desconto, pix_list
		FROM cobvs WHERE txid = $1`, txid)
	return r.scanRow(row)
}

func (r *CobVRepository) FindAll(filters interfaces.CobVFilters) ([]models.CobV, error) {
	query := `
		SELECT txid, status, chave, revisao, data_de_vencimento, validade_apos_vencimento,
		       criacao, valor_original, pix_copia_e_cola, solicitacao_pagador,
		       devedor, loc, info_adicionais, valor_multa, valor_juros,
		       valor_abatimento, valor_desconto, pix_list
		FROM cobvs WHERE 1=1`

	args := make([]any, 0)
	argIndex := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}
	if filters.DataDeVencimento != "" {
		query += fmt.Sprintf(" AND data_de_vencimento = $%d", argIndex)
		args = append(args, filters.DataDeVencimento)
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

func (r *CobVRepository) Update(cobv models.CobV) error {
	devedorJSON, err := json.Marshal(cobv.Devedor)
	if err != nil {
		return fmt.Errorf("marshal devedor: %w", err)
	}
	locJSON, err := marshalNullable(cobv.Loc)
	if err != nil {
		return fmt.Errorf("marshal loc: %w", err)
	}
	infoJSON, err := json.Marshal(cobv.InfoAdicionais)
	if err != nil {
		return fmt.Errorf("marshal info_adicionais: %w", err)
	}
	multaJSON, err := marshalNullable(cobv.Valor.Multa)
	if err != nil {
		return fmt.Errorf("marshal valor_multa: %w", err)
	}
	jurosJSON, err := marshalNullable(cobv.Valor.Juros)
	if err != nil {
		return fmt.Errorf("marshal valor_juros: %w", err)
	}
	abatimentoJSON, err := marshalNullable(cobv.Valor.Abatimento)
	if err != nil {
		return fmt.Errorf("marshal valor_abatimento: %w", err)
	}
	descontoJSON, err := marshalNullable(cobv.Valor.Desconto)
	if err != nil {
		return fmt.Errorf("marshal valor_desconto: %w", err)
	}
	pixListJSON, err := json.Marshal(cobv.Pix)
	if err != nil {
		return fmt.Errorf("marshal pix_list: %w", err)
	}

	result, err := r.db.Exec(`
		UPDATE cobvs SET
			status = $1, chave = $2, revisao = $3,
			data_de_vencimento = $4, validade_apos_vencimento = $5, criacao = $6,
			valor_original = $7, pix_copia_e_cola = $8, solicitacao_pagador = $9,
			devedor = $10, loc = $11, info_adicionais = $12,
			valor_multa = $13, valor_juros = $14, valor_abatimento = $15,
			valor_desconto = $16, pix_list = $17
		WHERE txid = $18`,
		cobv.Status, cobv.Chave, cobv.Revisao,
		cobv.Calendario.DataDeVencimento, cobv.Calendario.ValidadeAposVencimento, cobv.Calendario.Criacao,
		cobv.Valor.Original,
		nullableString(cobv.PixCopiaECola), nullableString(cobv.SolicitacaoPagador),
		devedorJSON, locJSON, infoJSON,
		multaJSON, jurosJSON, abatimentoJSON, descontoJSON,
		pixListJSON,
		cobv.TxID,
	)
	if err != nil {
		return err
	}
	return requireOneRowAffected(result, "cobrança com vencimento", cobv.TxID)
}

func (r *CobVRepository) scanRow(row *sql.Row) (*models.CobV, error) {
	var (
		txid, status, chave, dataVencimento, valorOriginal string
		revisao, validadeAposVencimento                    int
		pixCopiaECola, solPagador                          sql.NullString
		criacao                                            time.Time
		devedorJSON, locJSON, infoJSON                     []byte
		multaJSON, jurosJSON, abatimentoJSON               []byte
		descontoJSON, pixListJSON                          []byte
	)

	err := row.Scan(
		&txid, &status, &chave, &revisao, &dataVencimento, &validadeAposVencimento,
		&criacao, &valorOriginal, &pixCopiaECola, &solPagador,
		&devedorJSON, &locJSON, &infoJSON,
		&multaJSON, &jurosJSON, &abatimentoJSON, &descontoJSON,
		&pixListJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return buildCobV(
		txid, status, chave, revisao, dataVencimento, validadeAposVencimento,
		criacao, valorOriginal, pixCopiaECola, solPagador,
		devedorJSON, locJSON, infoJSON,
		multaJSON, jurosJSON, abatimentoJSON, descontoJSON,
		pixListJSON,
	)
}

func (r *CobVRepository) scanRows(rows *sql.Rows) ([]models.CobV, error) {
	result := make([]models.CobV, 0)
	for rows.Next() {
		var (
			txid, status, chave, dataVencimento, valorOriginal string
			revisao, validadeAposVencimento                    int
			pixCopiaECola, solPagador                          sql.NullString
			criacao                                            time.Time
			devedorJSON, locJSON, infoJSON                     []byte
			multaJSON, jurosJSON, abatimentoJSON               []byte
			descontoJSON, pixListJSON                          []byte
		)

		if err := rows.Scan(
			&txid, &status, &chave, &revisao, &dataVencimento, &validadeAposVencimento,
			&criacao, &valorOriginal, &pixCopiaECola, &solPagador,
			&devedorJSON, &locJSON, &infoJSON,
			&multaJSON, &jurosJSON, &abatimentoJSON, &descontoJSON,
			&pixListJSON,
		); err != nil {
			return nil, err
		}

		cobv, err := buildCobV(
			txid, status, chave, revisao, dataVencimento, validadeAposVencimento,
			criacao, valorOriginal, pixCopiaECola, solPagador,
			devedorJSON, locJSON, infoJSON,
			multaJSON, jurosJSON, abatimentoJSON, descontoJSON,
			pixListJSON,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, *cobv)
	}
	return result, rows.Err()
}

func buildCobV(
	txid, status, chave string,
	revisao int,
	dataVencimento string,
	validadeAposVencimento int,
	criacao time.Time,
	valorOriginal string,
	pixCopiaECola, solPagador sql.NullString,
	devedorJSON, locJSON, infoJSON []byte,
	multaJSON, jurosJSON, abatimentoJSON, descontoJSON []byte,
	pixListJSON []byte,
) (*models.CobV, error) {
	cobv := models.CobV{
		TxID:    txid,
		Status:  status,
		Chave:   chave,
		Revisao: revisao,
		Calendario: models.CalendarioCobV{
			DataDeVencimento:       dataVencimento,
			ValidadeAposVencimento: validadeAposVencimento,
			Criacao:                criacao,
		},
		Valor:              models.CobVValor{Original: valorOriginal},
		PixCopiaECola:      pixCopiaECola.String,
		SolicitacaoPagador: solPagador.String,
	}

	if len(devedorJSON) > 0 {
		if err := json.Unmarshal(devedorJSON, &cobv.Devedor); err != nil {
			return nil, fmt.Errorf("unmarshal devedor: %w", err)
		}
	}
	if len(locJSON) > 0 {
		var loc models.Loc
		if err := json.Unmarshal(locJSON, &loc); err != nil {
			return nil, fmt.Errorf("unmarshal loc: %w", err)
		}
		cobv.Loc = &loc
	}
	if len(infoJSON) > 0 {
		if err := json.Unmarshal(infoJSON, &cobv.InfoAdicionais); err != nil {
			return nil, fmt.Errorf("unmarshal info_adicionais: %w", err)
		}
	}
	if len(multaJSON) > 0 {
		var multa models.ValorComponente
		if err := json.Unmarshal(multaJSON, &multa); err != nil {
			return nil, fmt.Errorf("unmarshal valor_multa: %w", err)
		}
		cobv.Valor.Multa = &multa
	}
	if len(jurosJSON) > 0 {
		var juros models.ValorComponente
		if err := json.Unmarshal(jurosJSON, &juros); err != nil {
			return nil, fmt.Errorf("unmarshal valor_juros: %w", err)
		}
		cobv.Valor.Juros = &juros
	}
	if len(abatimentoJSON) > 0 {
		var abatimento models.ValorComponente
		if err := json.Unmarshal(abatimentoJSON, &abatimento); err != nil {
			return nil, fmt.Errorf("unmarshal valor_abatimento: %w", err)
		}
		cobv.Valor.Abatimento = &abatimento
	}
	if len(descontoJSON) > 0 {
		var desconto models.DescontoCobV
		if err := json.Unmarshal(descontoJSON, &desconto); err != nil {
			return nil, fmt.Errorf("unmarshal valor_desconto: %w", err)
		}
		cobv.Valor.Desconto = &desconto
	}
	if len(pixListJSON) > 0 {
		if err := json.Unmarshal(pixListJSON, &cobv.Pix); err != nil {
			return nil, fmt.Errorf("unmarshal pix_list: %w", err)
		}
	}

	return &cobv, nil
}
