package src

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type WorkerAbility int

const (
	C WorkerAbility = 0
	U WorkerAbility = 1
	R WorkerAbility = 2
	D WorkerAbility = 3
)

func runCurd(
	rq *require.Assertions,
	migrators []*Migrator,
	opNum int,
	abilities []WorkerAbility,
	ids []string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	ctx := context.Background()
	last := 0
	getMigrator := func() *Migrator {
		return migrators[last%len(migrators)]
	}

	for i := 0; i < opNum; i++ {
		ability := C
		index := 0
		id := ""
		migrator := getMigrator()

		if len(ids) > 0 {
			ability = abilities[rand.Int()%len(abilities)]
			index = rand.Intn(len(ids))
			id = ids[index]
		}

		switch ability {
		case C:
			m, err := migrator.Create(ctx, Model{Msg: "0"})
			rq.NoError(err)
			ids = append(ids, m.ID)
		case U:
			rq.NoError(migrator.Update(ctx, Model{ID: id, Msg: fmt.Sprintf("%d", last)}))
		case R:
			_, err := migrator.Read(ctx, id)
			rq.NoError(err)
		case D:
			err := migrator.Delete(ctx, id)
			rq.NoError(err)
			copy(ids[index:], ids[index+1:])
		}
	}
}

func RunCurdConcurrency(
	rq *require.Assertions,
	migrators []*Migrator,
	worker int,
	opNum int,
	abilities ...WorkerAbility,
) {
	var wg sync.WaitGroup

	wg.Add(worker)
	for i := 0; i < worker; i++ {
		go runCurd(rq, migrators, opNum, abilities, []string{}, &wg)
	}

	wg.Wait()
}

const sqlFetchInconsistentModel = `
	SELECT 
		origin.ID, origin.MSG, migrated.MSG
	FROM 
	   %s.models AS origin
	LEFT JOIN
		(
		    SELECT ID, MSG FROM %s.models
		    UNION ALL
		    SELECT ID, MSG FROM %s.models
		    UNION ALL
		    SELECT ID, MSG FROM %s.models
		) AS migrated
	ON 
		origin.ID = migrated.ID
	WHERE
	    origin.MSG != migrated.MSG
`

func ValidateConsistence(rq *require.Assertions, db *gorm.DB) {
	var models []struct {
		ID          string
		OriginMsg   string
		MigratedMsg string
	}

	sql := fmt.Sprintf(sqlFetchInconsistentModel, dbOrigin, dbMigrated1, dbMigrated2, dbMigrated3)
	result := db.Raw(sql).Scan(&models)
	rq.NoError(result.Error)

	if len(models) > 0 {
		var ret string
		for _, m := range models {
			ret += fmt.Sprintf("%s %s %s\r\n", m.ID, m.OriginMsg, m.MigratedMsg)
		}
		rq.Failf("inconsistency", "found %s", ret)
	}
}
