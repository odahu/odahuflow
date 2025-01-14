package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelTrainingTable = "odahu_operator_training"
)

var (
	log       = logf.Log.WithName("model-training--repository--postgres")
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
)

type TrainingRepo struct {
	DB *sql.DB
}

func (repo TrainingRepo) GetModelTraining(ctx context.Context, tx *sql.Tx, id string) (*training.ModelTraining, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	mt := new(training.ModelTraining)

	query, args, err := sq.
		Select("id", "spec", "status", "deletionmark", "created", "updated").
		From(ModelTrainingTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = qrr.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, odahuErrors.NotFoundError{Entity: id}
	case err != nil:
		log.Error(err, "error during sql query")
		return nil, err
	default:
		return mt, nil
	}

}

func (repo TrainingRepo) GetModelTrainingList(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]training.ModelTraining, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &FirstPage,
		Size:   &MaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	offset := *listOptions.Size * (*listOptions.Page)

	sb := sq.Select("id, spec, status, deletionmark, created, updated").From("odahu_operator_training").
		OrderBy("id").
		Offset(uint64(offset)).
		Limit(uint64(*listOptions.Size)).PlaceholderFormat(sq.Dollar)

	sb = utils.TransformFilter(sb, listOptions.Filter)
	stmt, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := qrr.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "error during rows.Close()")
		}
		if err := rows.Err(); err != nil {
			log.Error(err, "error during rows iterating")
		}
	}()

	mts := make([]training.ModelTraining, 0)

	for rows.Next() {
		mt := new(training.ModelTraining)
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		mts = append(mts, *mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mts, nil

}

func (repo TrainingRepo) DeleteModelTraining(ctx context.Context, tx *sql.Tx, id string) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(ModelTrainingTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: id}
	}

	return nil

}

func (repo TrainingRepo) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	return utils.SetDeletionMark(ctx, qrr, ModelTrainingTable, id, value)
}

func (repo TrainingRepo) UpdateModelTraining(ctx context.Context, tx *sql.Tx, mt *training.ModelTraining) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(ModelTrainingTable).
		Set("spec", mt.Spec).
		Set("status", mt.Status).
		Set("updated", mt.UpdatedAt).
		Where(sq.Eq{"id": mt.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: mt.ID}
	}

	return nil
}

func (repo TrainingRepo) UpdateModelTrainingStatus(
	ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelTrainingStatus) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(ModelTrainingTable).
		Set("status", s).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: id}
	}

	return nil
}

func (repo TrainingRepo) SaveModelTraining(ctx context.Context, tx *sql.Tx, mt *training.ModelTraining) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(ModelTrainingTable).
		Columns("id", "spec", "status", "created", "updated").
		Values(mt.ID, mt.Spec, mt.Status, mt.CreatedAt, mt.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = qrr.ExecContext(
		ctx,
		stmt,
		args...,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mt.ID}
		}
		return err
	}
	return nil

}

func (repo TrainingRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return repo.DB.BeginTx(ctx, txOptions)
}
