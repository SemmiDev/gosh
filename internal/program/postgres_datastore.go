package program

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProgramDataStore struct {
	pool *pgxpool.Pool
}

func NewProgramDataStore(pool *pgxpool.Pool) *ProgramDataStore {
	return &ProgramDataStore{
		pool: pool,
	}
}

func (p *ProgramDataStore) CreateProgram(ctx context.Context, program Program) (Program, error) {
	var id int64
	err := p.pool.QueryRow(ctx, "INSERT INTO program (name, description) VALUES ($1, $2) RETURNING id", program.Name, program.Description).Scan(&id)
	if err != nil {
		return Program{}, err
	}

	program.ID = id
	return program, nil
}

func (p *ProgramDataStore) GetProgramDetails(ctx context.Context, id int64) (Program, error) {
	var program Program
	err := p.pool.QueryRow(ctx, "SELECT id, name, description FROM program WHERE id = $1", id).Scan(&program.ID, &program.Name, &program.Description)
	return program, err
}

func (p *ProgramDataStore) SearchTerm(ctx context.Context, q string) ([]Program, error) {
	// english, indonesia, etc..
	rows, err := p.pool.Query(ctx, "SELECT id, name, description FROM program WHERE ts @@ to_tsquery('english', $1)", q)
	if err != nil {
		return nil, err
	}

	var programs []Program
	for rows.Next() {
		var program Program
		err = rows.Scan(&program.ID, &program.Name, &program.Description)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	return programs, nil
}
