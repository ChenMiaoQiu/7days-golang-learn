package session

import "github.com/ChenMiaoQiu/7days-golang-learn/ORM/day6/geeORM/log"

func (s *Session) Begin() (err error) {
	log.Info("transation begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transation commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) RollBack() (err error) {
	log.Info("transation rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	return
}
