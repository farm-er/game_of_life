package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)




type Border struct {

	Height int

	Width int

}




type Screen struct {

	// height of the simulation screen 
	Height int	

	// width of the simulation screen
	Width int

  S tcell.Screen

	Cells [][]bool

	NextCells [][]bool

	FrameDur time.Duration

	Fps int

	StartTime time.Time

	FrameCount int

	ContentX int

	Border *Border

}

func (s *Screen) Start() {
	
	defer s.S.Fini()

	s.StartTime = time.Now()
		

	eventChan := make(chan tcell.Event, 1)
	
	go func() {
		for {
		ev := s.S.PollEvent()
		eventChan <- ev
		}
	}()

	s.DrawBoard()

	for {

		select {
		case ev := <-eventChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEsc:
					s.S.Fini()
					os.Exit(0)
				}
			}
		default:
			// continue 
		}

		start := time.Now()


		s.DrawCells()

				// add the frame
		s.FrameCount++


		// make the status bar
		barStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

		barContent := fmt.Sprintf("cell 1 size(%v, %v) FPS: %.2f", s.Width, s.Height, float64(s.FrameCount)/time.Since(s.StartTime).Seconds())

		s.S.SetContent( 0, s.Height, ' ', []rune(barContent), barStyle)


		end := time.Since(start)

		if end < s.FrameDur {
			time.Sleep(s.FrameDur-end-7)
		}

		s.S.Show()
	}

}


func NewScreen() ( *Screen, error) {

	s, e := tcell.NewScreen()

	if e != nil {
		return nil, e
	}

	s.Init()

	w, h := s.Size()	

	
	cells := make([][]bool, h-3)
	nextCells := make([][]bool, h-3)
	for i:= range cells {
		cells[i] = make([]bool, w-2)
		nextCells[i] = make([]bool, w-2)
	}

	nextCells[5][5] = true
	nextCells[5][4] = true
	nextCells[5][3] = true

	nextCells[20][20] = true
	nextCells[20][19] = true
	nextCells[20][18] = true


	nextCells[5][10] = true
	nextCells[5][11] = true
	nextCells[5][12] = true
	nextCells[4][9] = true
	nextCells[4][10] = true
	nextCells[4][11] = true
	


	return &Screen{
		Height: h-1,
		Width: w,
		S: s,
		Cells: cells,
		NextCells: nextCells,
		Fps: 60,
		FrameDur: time.Second/60,
		FrameCount: 0,
		ContentX: h,
		Border: &Border{
			Height: h-3,
			Width: w-2,
		},
	}, nil

}

func (s *Screen) DrawBoard() {

	for j:= 0;j<s.Height-1; j++ {
		s.S.SetContent( 0, j, '│', nil, tcell.Style{})
		s.S.SetContent( s.Width-1, j, '│', nil, tcell.Style{})
	}


	for j:= 0;j<s.Width; j++ {
		s.S.SetContent( j, s.Height-1, '─', nil, tcell.Style{})
		s.S.SetContent( j, 0, '─', nil, tcell.Style{})
	}


	s.S.SetContent( 0, 0, '┌', nil, tcell.Style{})
	s.S.SetContent( s.Width-1, 0, '┐', nil, tcell.Style{})
	s.S.SetContent( 0, s.Height-1, '└', nil, tcell.Style{})
	s.S.SetContent( s.Width-1, s.Height-1, '┘', nil, tcell.Style{})

}


func (s *Screen) DrawCells() {

	value := ' '
	
	for i:= range s.Cells {
		copy(s.Cells[i], s.NextCells[i])
	}
	
	for i:= range s.Cells {
		for j:= range s.Cells[i] {
			// here do the trick
			s.UpdateCellState( i, j)

			if s.NextCells[i][j] {
				s.S.SetContent( j+1, i+1, value, nil, tcell.StyleDefault.Background(tcell.ColorWhite))
			}else {
				s.S.SetContent( j+1, i+1, value, nil, tcell.Style{})
			}

		}
	}
}


func (s *Screen) UpdateCellState( i, j int) {
		
	counter := s.aliveCells(i, j)

	if ( s.Cells[i][j] && ( counter < 2 || counter > 3)) { s.NextCells[i][j] = false; return}
	if ( !s.Cells[i][j] && ( counter == 3 ) ) { s.NextCells[i][j] = true; return }

	s.NextCells[i][j] = s.Cells[i][j]

}

func (s *Screen) aliveCells( i, j int	) int {

	counter := 0

	h := s.Border.Height - 1
	w := s.Border.Width - 1


	if ( j < w && s.Cells[i][j+1] ) {counter++}
	if ( j > 0 && s.Cells[i][j-1] ) {counter++}
	if ( i < h && j < w && s.Cells[i+1][j+1]) {counter++}
	if ( i < h && j > 0 && s.Cells[i+1][j-1] ) {counter++}
	if ( i < h && s.Cells[i+1][j] ) {counter++}
	if ( i > 0 && s.Cells[i-1][j] ) {counter++}
	if ( i > 0 && j < w  && s.Cells[i-1][j+1] ) {counter++}
	if ( i > 0 && j > 0 && s.Cells[i-1][j-1] ) {counter++}

	return counter
}

