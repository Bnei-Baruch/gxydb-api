package domain

import (
	"context"
	"database/sql"

	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/sqlutil"
)

type RoomStatisticsManager struct {
	db common.DBInterface
}

func NewRoomStatisticsManager(db common.DBInterface) *RoomStatisticsManager {
	return &RoomStatisticsManager{
		db: db,
	}
}

func (m *RoomStatisticsManager) OnAir(roomID string) error {
	roomStats, err := m.getOrCreate(roomID)
	if err != nil {
		return err
	}
	roomStats.OnAir++
	_, err = roomStats.Update(m.db, boil.Infer())
	return err
}

func (m *RoomStatisticsManager) GetAll() ([]*models.RoomStatistic, error) {
	// Note: Room relationship no longer exists (room_id is string gateway_uid, no FK constraint)
	// Load rooms manually if needed
	return models.RoomStatistics().All(m.db)
}

func (m *RoomStatisticsManager) Reset(ctx context.Context) error {
	return sqlutil.InTx(ctx, m.db, func(tx *sql.Tx) error {
		rowsAff, err := models.RoomStatistics().DeleteAll(tx)
		if err != nil {
			return pkgerr.WithStack(err)
		}

		log.Ctx(ctx).Info().Int64("deleted", rowsAff).Msg("delete rooms statistics")

		return nil
	})
}

func (m *RoomStatisticsManager) getOrCreate(roomID string) (*models.RoomStatistic, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	roomStats, err := models.FindRoomStatistic(tx, roomID)
	if err == nil {
		return roomStats, tx.Commit()
	}

	roomStats = &models.RoomStatistic{
		RoomID: roomID,
		OnAir:  0,
	}
	if err := roomStats.Insert(tx, boil.Infer()); err != nil {
		return nil, err
	}

	return roomStats, tx.Commit()
}

func (m *RoomStatisticsManager) update(roomStats *models.RoomStatistic) error {
	return sqlutil.InTx(context.TODO(), m.db, func(tx *sql.Tx) error {
		_, err := roomStats.Update(tx, boil.Infer())
		if err != nil {
			return pkgerr.WithStack(err)
		}
		return nil
	})
}
