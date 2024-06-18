package oviewer

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestRoot_moveVertical(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()
	type fields struct {
		fileName string
		topLN    int
	}
	type want struct {
		top    int
		bottom int
		pgup   int
		pgdn   int
		pghfup int
		pghfdn int
		upOne  int
		dnOne  int
		up     int
		dn     int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test",
			fields: fields{
				fileName: filepath.Join(testdata, "test3.txt"),
				topLN:    1000,
			},
			want: want{
				top:    0,
				bottom: 12322,
				pgup:   976,
				pgdn:   1024,
				pghfup: 988,
				pghfdn: 1012,
				upOne:  999,
				dnOne:  1001,
				up:     996,
				dn:     1004,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := Open(tt.fields.fileName)
			if err != nil {
				t.Fatal(err)
			}
			root.prepareScreen()
			ctx := context.Background()
			root.everyUpdate(context.Background())
			for !root.Doc.BufEOF() {
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveTop(ctx)
			if root.Doc.topLN != tt.want.top {
				t.Errorf("top got %d, want %d", root.Doc.topLN, tt.want.top)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveBottom(ctx)
			if root.Doc.topLN != tt.want.bottom {
				t.Errorf("botoom got %d, want %d", root.Doc.topLN, tt.want.bottom)
			}
			root.Doc.topLN = tt.fields.topLN
			root.movePgUp(ctx)
			if root.Doc.topLN != tt.want.pgup {
				t.Errorf("pgup got %d, want %d", root.Doc.topLN, tt.want.pgup)
			}
			root.Doc.topLN = tt.fields.topLN
			root.movePgDn(ctx)
			if root.Doc.topLN != tt.want.pgdn {
				t.Errorf("pgdown got %d, want %d", root.Doc.topLN, tt.want.pgdn)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveHfUp(ctx)
			if root.Doc.topLN != tt.want.pghfup {
				t.Errorf("pghfup got %d, want %d", root.Doc.topLN, tt.want.pghfup)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveHfDn(ctx)
			if root.Doc.topLN != tt.want.pghfdn {
				t.Errorf("pghfdn got %d, want %d", root.Doc.topLN, tt.want.pghfdn)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveUpOne(ctx)
			if root.Doc.topLN != tt.want.upOne {
				t.Errorf("upOne got %d, want %d", root.Doc.topLN, tt.want.upOne)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveDownOne(ctx)
			if root.Doc.topLN != tt.want.dnOne {
				t.Errorf("dnOne got %d, want %d", root.Doc.topLN, tt.want.dnOne)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveUp(4)
			if root.Doc.topLN != tt.want.up {
				t.Errorf("up got %d, want %d", root.Doc.topLN, tt.want.up)
			}
			root.Doc.topLN = tt.fields.topLN
			root.moveDown(4)
			if root.Doc.topLN != tt.want.dn {
				t.Errorf("dn got %d, want %d", root.Doc.topLN, tt.want.dn)
			}
		})
	}
}

func TestRoot_moveSection(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()
	type fields struct {
		fileName         string
		topLN            int
		sectionDelimiter string
	}
	type want struct {
		nextSection int
		prevSection int
		lastSection int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test no section",
			fields: fields{
				fileName: filepath.Join(testdata, "test3.txt"),
				topLN:    1000,
			},
			want: want{
				nextSection: 1024,
				prevSection: 976,
				lastSection: 1000,
			},
		},
		{
			name: "Test section",
			fields: fields{
				fileName:         filepath.Join(testdata, "test3.txt"),
				topLN:            1000,
				sectionDelimiter: "000",
			},
			want: want{
				nextSection: 1999,
				prevSection: 999,
				lastSection: 11999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := Open(tt.fields.fileName)
			if err != nil {
				t.Fatal(err)
			}
			root.prepareScreen()
			ctx := context.Background()
			root.everyUpdate(context.Background())
			for !root.Doc.BufEOF() {
			}
			root.Doc.topLN = tt.fields.topLN
			root.setSectionDelimiter(tt.fields.sectionDelimiter)
			root.nextSection(ctx)
			if root.Doc.topLN != tt.want.nextSection {
				t.Errorf("nextSection got %d, want %d", root.Doc.topLN, tt.want.nextSection)
			}
			root.Doc.topLN = tt.fields.topLN
			root.prevSection(ctx)
			if root.Doc.topLN != tt.want.prevSection {
				t.Errorf("prevSection got %d, want %d", root.Doc.topLN, tt.want.prevSection)
			}

			root.Doc.topLN = tt.fields.topLN
			root.lastSection(ctx)
			if root.Doc.topLN != tt.want.lastSection {
				t.Errorf("lastSection got %d, want %d", root.Doc.topLN, tt.want.lastSection)
			}
			root.nextSection(ctx)
			if err := root.Doc.moveNextSection(ctx); err == nil {
				t.Errorf("no more section? %d", root.Doc.topLN)
			}
			root.moveTop(ctx)
			root.prevSection(ctx)
			if err := root.Doc.movePrevSection(ctx); err == nil {
				t.Errorf("no more section? %d", root.Doc.topLN)
			}
		})
	}
}

func TestRoot_moveLateral(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()
	type fields struct {
		fileName string
		x        int
	}
	type want struct {
		leftOne    int
		rightOne   int
		hfLeft     int
		hfRight    int
		widthLeft  int
		widthRight int
		bgLeft     int
		endRight   int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test",
			fields: fields{
				fileName: filepath.Join(testdata, "normal.txt"),
				x:        10,
			},
			want: want{
				leftOne:    9,
				rightOne:   11,
				hfLeft:     0,
				hfRight:    50,
				widthLeft:  9,
				widthRight: 11,
				bgLeft:     0,
				endRight:   20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := Open(tt.fields.fileName)
			if err != nil {
				t.Fatal(err)
			}
			root.prepareScreen()
			ctx := context.Background()
			root.everyUpdate(context.Background())
			for !root.Doc.BufEOF() {
			}
			root.Doc.ColumnMode = false
			root.Doc.WrapMode = false
			root.Doc.x = tt.fields.x
			root.moveLeftOne(ctx)
			if root.Doc.x != tt.want.leftOne {
				t.Errorf("leftOne got %d, want %d", root.Doc.x, tt.want.leftOne)
			}
			root.Doc.x = tt.fields.x
			root.moveRightOne(ctx)
			if root.Doc.x != tt.want.rightOne {
				t.Errorf("rightOne got %d, want %d", root.Doc.x, tt.want.rightOne)
			}
			root.Doc.x = tt.fields.x
			root.moveHfLeft(ctx)
			if root.Doc.x != tt.want.hfLeft {
				t.Errorf("moveHfLeft got %d, want %d", root.Doc.x, tt.want.hfLeft)
			}
			root.Doc.x = tt.fields.x
			root.moveHfRight(ctx)
			if root.Doc.x != tt.want.hfRight {
				t.Errorf("moveHfRight got %d, want %d", root.Doc.x, tt.want.hfRight)
			}
			root.Doc.x = tt.fields.x
			root.moveWidthLeft(ctx)
			if root.Doc.x != tt.want.widthLeft {
				t.Errorf("moveWidthLeft got %d, want %d", root.Doc.x, tt.want.widthLeft)
			}
			root.Doc.x = tt.fields.x
			root.moveWidthRight(ctx)
			if root.Doc.x != tt.want.widthRight {
				t.Errorf("moveWidthRight got %d, want %d", root.Doc.x, tt.want.widthRight)
			}
			root.Doc.x = tt.fields.x
			root.moveBeginLeft(ctx)
			if root.Doc.x != tt.want.bgLeft {
				t.Errorf("moveBeginLeft got %d, want %d", root.Doc.x, tt.want.bgLeft)
			}
			root.Doc.x = tt.fields.x
			root.moveEndRight(ctx)
			if root.Doc.x != tt.want.endRight {
				t.Errorf("moveEndRight got %d, want %d", root.Doc.x, tt.want.endRight)
			}
		})
	}
}

func TestRoot_moveColumn(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()
	type fields struct {
		fileName     string
		columnCursor int
	}
	type want struct {
		leftOne  int
		rightOne int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test no section",
			fields: fields{
				fileName:     filepath.Join(testdata, "MOCK_DATA.csv"),
				columnCursor: 4,
			},
			want: want{
				leftOne:  3,
				rightOne: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := Open(tt.fields.fileName)
			if err != nil {
				t.Fatal(err)
			}
			root.prepareScreen()
			ctx := context.Background()
			root.everyUpdate(context.Background())
			for !root.Doc.BufEOF() {
			}
			root.Doc.ColumnMode = true
			root.Doc.setDelimiter(",")

			root.Doc.columnCursor = tt.fields.columnCursor
			root.moveLeftOne(ctx)
			if root.Doc.columnCursor != tt.want.leftOne {
				t.Errorf("leftOne got %d, want %d", root.Doc.columnCursor, tt.want.leftOne)
			}
			root.Doc.columnCursor = tt.fields.columnCursor
			root.moveRightOne(ctx)
			if root.Doc.columnCursor != tt.want.rightOne {
				t.Errorf("rightOne got %d, want %d", root.Doc.columnCursor, tt.want.rightOne)
			}
		})
	}
}
