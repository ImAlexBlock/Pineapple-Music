package service

import (
	"fmt"

	"pineapple-music/internal/model"

	"gorm.io/gorm"
)

type PlaylistService struct {
	DB *gorm.DB
}

func (s *PlaylistService) Create(name string) (*model.Playlist, error) {
	pl := &model.Playlist{Name: name}
	if err := s.DB.Create(pl).Error; err != nil {
		return nil, err
	}
	return pl, nil
}

func (s *PlaylistService) List() ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := s.DB.Find(&playlists).Error
	return playlists, err
}

func (s *PlaylistService) Get(id uint) (*model.Playlist, error) {
	var pl model.Playlist
	if err := s.DB.Preload("Tracks", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC").Preload("Track")
	}).First(&pl, id).Error; err != nil {
		return nil, err
	}
	return &pl, nil
}

func (s *PlaylistService) Delete(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("playlist_id = ?", id).Delete(&model.PlaylistTrack{})
		return tx.Delete(&model.Playlist{}, id).Error
	})
}

func (s *PlaylistService) AddTrack(playlistID, trackID uint) error {
	// Get max position
	var maxPos int
	s.DB.Model(&model.PlaylistTrack{}).
		Where("playlist_id = ?", playlistID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPos)

	pt := &model.PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    trackID,
		Position:   maxPos + 1,
	}
	return s.DB.Create(pt).Error
}

func (s *PlaylistService) RemoveTrack(playlistID, trackID uint) error {
	return s.DB.Where("playlist_id = ? AND track_id = ?", playlistID, trackID).
		Delete(&model.PlaylistTrack{}).Error
}

func (s *PlaylistService) ReorderTracks(playlistID uint, trackIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for i, tid := range trackIDs {
			result := tx.Model(&model.PlaylistTrack{}).
				Where("playlist_id = ? AND track_id = ?", playlistID, tid).
				Update("position", i+1)
			if result.RowsAffected == 0 {
				return fmt.Errorf("track %d not in playlist", tid)
			}
		}
		return nil
	})
}
