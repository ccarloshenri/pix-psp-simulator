package sqlrepo

import "database/sql"

// RunMigrations creates all required tables if they do not yet exist.
// It is idempotent and safe to call on every startup.
func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cobs (
			txid                  VARCHAR(35)  PRIMARY KEY,
			status                VARCHAR(50)  NOT NULL,
			chave                 VARCHAR(77)  NOT NULL,
			revisao               INTEGER      NOT NULL DEFAULT 0,
			expiracao             INTEGER      NOT NULL,
			criacao               TIMESTAMPTZ  NOT NULL,
			valor_original        VARCHAR(20)  NOT NULL,
			modalidade_alteracao  INTEGER      DEFAULT 0,
			location              VARCHAR(77),
			pix_copia_e_cola      VARCHAR(512),
			solicitacao_pagador   VARCHAR(140),
			devedor               JSONB,
			loc                   JSONB,
			info_adicionais       JSONB        DEFAULT '[]',
			retirada              JSONB,
			pix_list              JSONB        DEFAULT '[]'
		);

		CREATE TABLE IF NOT EXISTS cobvs (
			txid                     VARCHAR(35)  PRIMARY KEY,
			status                   VARCHAR(50)  NOT NULL,
			chave                    VARCHAR(77)  NOT NULL,
			revisao                  INTEGER      NOT NULL DEFAULT 0,
			data_de_vencimento       VARCHAR(10)  NOT NULL,
			validade_apos_vencimento INTEGER      DEFAULT 30,
			criacao                  TIMESTAMPTZ  NOT NULL,
			valor_original           VARCHAR(20)  NOT NULL,
			pix_copia_e_cola         VARCHAR(512),
			solicitacao_pagador      VARCHAR(140),
			devedor                  JSONB        NOT NULL DEFAULT '{}',
			loc                      JSONB,
			info_adicionais          JSONB        DEFAULT '[]',
			valor_multa              JSONB,
			valor_juros              JSONB,
			valor_abatimento         JSONB,
			valor_desconto           JSONB,
			pix_list                 JSONB        DEFAULT '[]'
		);

		CREATE TABLE IF NOT EXISTS pix_payments (
			end_to_end_id  VARCHAR(50)  PRIMARY KEY,
			txid           VARCHAR(35)  NOT NULL,
			valor          VARCHAR(20)  NOT NULL,
			horario        TIMESTAMPTZ  NOT NULL,
			infopagador    VARCHAR(200),
			devolucoes     JSONB        DEFAULT '[]'
		);
	`)
	return err
}
