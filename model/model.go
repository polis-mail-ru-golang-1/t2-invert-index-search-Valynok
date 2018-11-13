package model

import (
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

type Model struct {
	db *pg.DB
}

func New(db *pg.DB) Model {
	return Model{
		db: db,
	}
}

type File struct {
	Id   int
	Name string `sql:"name"`
}

func (m Model) ClearModel() {
	_, err := m.db.Exec("TRUNCATE counters CASCADE")
	if err != nil {
		Logger.Error(err)
	}

	_, err = m.db.Exec("TRUNCATE files CASCADE")
	if err != nil {
		Logger.Error(err)
	}

	_, err = m.db.Exec("TRUNCATE words CASCADE")
	if err != nil {
		Logger.Error(err)
	}
}

type Word struct {
	Id   int
	Word string `sql:"word"`
}

type Counters struct {
	Id      int
	WordId  int `sql:"wordId"`
	FileId  int `sql:"fileId"`
	Counter int `sql:"counter"`
}

func (m Model) AddCounters(wordid int, fileid int, counter int) {
	counters := Counters{
		WordId:  wordid,
		FileId:  fileid,
		Counter: counter,
	}

	Logger.Debug(counters)

	m.db.Insert(&counters)
}

func (m Model) AddCountersBulk(counters []Counters) {
	m.db.Insert(&counters)
}

func (m Model) AddCoutersBulk(counters []Counters) {
	err := m.db.Insert(&counters)
	if err != nil {
		Logger.Error(err)
	}
}

func (m Model) GetOrAddWord(wordname string) Word {
	word := Word{
		Word: wordname,
	}

	_, err := m.db.Model(&word).
		Column("id").
		Where("word = ?word").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()

	if err != nil {
		Logger.Error(err)
	}

	return word
}

func (m Model) AddWordBulk(words []string) []Word {
	wordsToAdd := make([]Word, len(words))

	for i, val := range words {
		wordsToAdd[i] = Word{
			Word: val,
		}
	}

	err := m.db.Insert(&wordsToAdd)
	if err != nil {
		Logger.Error(err)
	}

	return wordsToAdd
}

func (m Model) GetWord(wordName string) *Word {
	word := new(Word)
	m.db.Model(word).Where("word = ?", wordName).Select()
	return word
}

func (m Model) GetWords(wordNames []string) *[]Word {
	result := new([]Word)
	err := m.db.Model(result).WhereIn("word in (?)", pg.In(wordNames)).Select()

	if err != nil {
		Logger.Error(err)
	}

	return result
}

func (m Model) GetOrAddFile(filename string) File {
	file := File{
		Name: filename,
	}
	_, err := m.db.Model(&file).
		Column("id").
		Where("name = ?name").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()

	if err != nil {
		Logger.Error(err)
	}

	return file
}

func (m Model) GetFile(fileName string) *File {
	file := new(File)
	m.db.Model(file).Where("name = ?", fileName).Select()
	return file
}