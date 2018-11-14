package model

import (
	"github.com/go-pg/pg"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
	"go.uber.org/zap"
)

type Model struct {
	db *pg.DB
	l  *zap.SugaredLogger
}

func New(db *pg.DB, logger *zap.SugaredLogger) Model {
	return Model{
		db: db,
		l:  logger,
	}
}

type File struct {
	Id   int
	Name string `sql:"name"`
}

func (m Model) ClearModel() {
	_, err := m.db.Exec("TRUNCATE counters CASCADE")
	if err != nil {
		m.l.Error(err)
	}

	_, err = m.db.Exec("TRUNCATE files CASCADE")
	if err != nil {
		m.l.Error(err)
	}

	_, err = m.db.Exec("TRUNCATE words CASCADE")
	if err != nil {
		m.l.Error(err)
	}
}

type Word struct {
	Id   int
	Word string `sql:"word"`
}

type Counters struct {
	Id      int
	WordId  int `sql:"wordid"`
	FileId  int `sql:"fileid"`
	Counter int `sql:"counter"`
}

func (m Model) AddCounters(wordid int, fileid int, counter int) {
	counters := Counters{
		WordId:  wordid,
		FileId:  fileid,
		Counter: counter,
	}

	m.db.Insert(&counters)
}

func (m Model) AddCountersBulk(counters []Counters) {
	m.db.Insert(&counters)
}

func (m Model) AddCoutersBulk(counters []Counters) {
	err := m.db.Insert(&counters)
	if err != nil {
		m.l.Error(err)
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
		m.l.Error(err)
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
		m.l.Error(err)
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
		m.l.Error(err)
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
		m.l.Error(err)
	}

	return file
}

func (m Model) GetFile(fileName string) *File {
	file := new(File)
	m.db.Model(file).Where("name = ?", fileName).Select()
	return file
}

func (m Model) GetFiles(ids []int) []File {
	res := make([]File, 0)
	m.db.Model(&res).Where("id in (?)", pg.In(ids)).Select()

	return res
}

type CounterResult struct {
	File    string
	Counter int
}

func (m Model) getCounters(wordsIds []int) []Counters {

	var res []Counters
	err := m.db.Model(&Counters{}).
		Column("fileid").
		ColumnExpr("SUM(counter) as counter").
		WhereIn("wordid in (?)", pg.In(wordsIds)).
		Group("fileid").
		Select(&res)

	if err != nil {
		m.l.Error(err)
	}

	return res
}

func (m Model) GetCountersResult(words []string) []CounterResult {
	wordsDb := m.GetWords(words)
	wordsIds := make([]int, 0)
	for _, val := range *wordsDb {
		wordsIds = append(wordsIds, val.Id)
	}

	counters := m.getCounters(wordsIds)

	fileIds := make([]int, 0)
	for _, val := range counters {
		fileIds = append(fileIds, val.FileId)
	}

	files := m.GetFiles(fileIds)
	res := make([]CounterResult, 0, len(counters))
	for _, val := range counters {
		index := mapUtils.Find(len(files), func(i int) bool {
			return files[i].Id == val.FileId
		})

		fileName := files[index].Name

		res = append(res, CounterResult{File: fileName, Counter: val.Counter})
	}

	return res
}
