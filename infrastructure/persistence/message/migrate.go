package message

import (
	"context"
	"fmt"
)

func (p *persistence) migrate(ctx context.Context) error {
	if err := p.postgreSQL.Exec(ctx, `
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
				CREATE TYPE message_status AS ENUM ('PENDING', 'SENT');
			END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			content TEXT NOT NULL,
			phone VARCHAR(255) NOT NULL,
			status message_status NOT NULL
		);

		ALTER TABLE messages REPLICA IDENTITY FULL;
	`); err != nil {
		return fmt.Errorf("persistence.postgreSQL.Exec(): %w", err)
	}

	return nil
}
