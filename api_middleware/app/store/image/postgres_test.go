package image

import (
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPgStore_CheckBarcode(t *testing.T) {
	srv := preparePgStore(t)

	tx, err := srv.connPool.Begin()
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO images(bar_code, img_type, mime, local_filename) "+
		"VALUES ($1, $2, $3, $4)", "foobarblah", ImgTypeBackground, "png", "file.png")
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err, "failed to commit insert transaction")

	ok, err := srv.CheckBarcode("foobarblah123")
	require.NoError(t, err)
	assert.False(t, ok)

	ok, err = srv.CheckBarcode("foobarblah")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestPgStore_GetBackgroundIds(t *testing.T) {
	srv := preparePgStore(t)
	imgs := []Image{
		{
			Barcode:       "blah1",
			ImgType:       ImgTypeSrc,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah2",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah3",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah4",
			ImgType:       ImgTypeBackground,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah5",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah6",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah7",
			ImgType:       ImgTypeBackground,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah8",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah9",
			ImgType:       ImgTypeCommitted,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
	}

	for i := range imgs {
		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO images(bar_code, img_type, mime, local_filename) "+
			"VALUES ($1, $2, $3, $4) "+
			"RETURNING id",
			imgs[i].Barcode,
			imgs[i].ImgType,
			imgs[i].Mime,
			imgs[i].LocalFilename,
		)

		err = row.Scan(&imgs[i].ID)
		require.NoError(t, err, "failed to scan inserting img id")

		err = tx.Commit()
		require.NoError(t, err)
	}

	imgIds, err := srv.GetBackgroundIds()
	require.NoError(t, err)

	ok := reflect.DeepEqual([]int{imgs[3].ID, imgs[6].ID}, imgIds)
	assert.True(t, ok)

}

func TestPgStore_getImgByBarcode(t *testing.T) {
	srv := preparePgStore(t)

	tx, err := srv.connPool.Begin()
	require.NoError(t, err)

	row := tx.QueryRow("INSERT INTO images(bar_code, img_type, mime, local_filename) "+
		"VALUES ($1, $2, $3, $4) "+
		"RETURNING id",
		"foobarblah",
		ImgTypeBackground,
		"png",
		"gopher.png",
	)
	var id int
	err = row.Scan(&id)
	require.NoError(t, err, "failed to scan inserting img id")

	err = tx.Commit()
	require.NoError(t, err, "failed to commit insert transaction")

	img, err := srv.getImgByBarcode("foobarblah")
	require.NoError(t, err)

	assert.Equal(t, id, img.ID)
	assert.Equal(t, ImgTypeBackground, img.ImgType)
	assert.Equal(t, "png", img.Mime)
	assert.Equal(t, "gopher.png", img.LocalFilename)
}

func TestPgStore_putImage(t *testing.T) {
	srv := preparePgStore(t)
	img := Image{
		Barcode:       "foobarblahblah",
		ImgType:       ImgTypeSrc,
		Mime:          "png",
		LocalFilename: "gophergo.png",
	}

	id, err := srv.putImage(img)
	require.NoError(t, err)

	row := srv.connPool.QueryRow("SELECT bar_code, img_type, mime, local_filename "+
		"FROM images WHERE id = $1", id)

	var bc, imgtype, mime, localfname string
	err = row.Scan(&bc, &imgtype, &mime, &localfname)
	require.NoError(t, err)

	assert.Equal(t, img.Barcode, bc)
	assert.Equal(t, img.ImgType, imgtype)
	assert.Equal(t, img.Mime, mime)
	assert.Equal(t, img.LocalFilename, localfname)

}

func TestPgStore_getImage(t *testing.T) {
	srv := preparePgStore(t)
	imgs := []Image{
		{
			Barcode:       "blah1",
			ImgType:       ImgTypeSrc,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah2",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah3",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah4",
			ImgType:       ImgTypeBackground,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah5",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah6",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah7",
			ImgType:       ImgTypeBackground,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah8",
			ImgType:       ImgTypeDerived,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
		{
			Barcode:       "blah9",
			ImgType:       ImgTypeCommitted,
			Mime:          "png",
			LocalFilename: "gopher.png",
		},
	}

	for i := range imgs {
		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO images(bar_code, img_type, mime, local_filename) "+
			"VALUES ($1, $2, $3, $4) "+
			"RETURNING id",
			imgs[i].Barcode,
			imgs[i].ImgType,
			imgs[i].Mime,
			imgs[i].LocalFilename,
		)

		err = row.Scan(&imgs[i].ID)
		require.NoError(t, err, "failed to scan inserting img id")

		err = tx.Commit()
		require.NoError(t, err)
	}

	img, err := srv.getImage(imgs[3].ID)
	require.NoError(t, err)

	assert.Equal(t, imgs[3].Barcode, img.Barcode)
	assert.Equal(t, imgs[3].ImgType, img.ImgType)
	assert.Equal(t, imgs[3].Mime, img.Mime)
	assert.Equal(t, imgs[3].LocalFilename, img.LocalFilename)
}

func preparePgStore(t *testing.T) *PgStore {
	connStr := os.Getenv("DB_TEST")

	connConf, err := pgx.ParseConnectionString(connStr)

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 2,
		AfterConnect:   nil,
		AcquireTimeout: 60 * time.Millisecond,
	})
	require.NoError(t, err)

	st := PgStore{
		ConnStr:  connStr,
		connPool: connPool,
	}

	require.NoError(t, err)

	cleanupStorage(t, st.connPool)
	t.Cleanup(func() {
		cleanupStorage(t, st.connPool)
	})

	return &st
}

func cleanupStorage(t *testing.T, pool *pgx.ConnPool) {
	tx, err := pool.Begin()
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE images CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE users CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`ALTER SEQUENCE images_id_seq RESTART WITH 1`)
	require.NoError(t, err)
	_, err = tx.Exec(`ALTER SEQUENCE users_id_seq RESTART WITH 1`)
	err = tx.Commit()
	require.NoError(t, err)
}
