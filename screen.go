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

	CellsAlive int
	
	PausedTime time.Duration

}

func (s *Screen) PauseGame( eventChan chan tcell.Event) {

	start := time.Now()

	for {
		select {				
		case ps := <- eventChan:
			switch ps := ps.(type) {
				case *tcell.EventKey:
				switch ps.Key() {
					case tcell.KeyEsc:
						s.S.Fini()
						os.Exit(0)
					case tcell.KeyRune:
						if ps.Rune() == ' ' {
							s.PausedTime += time.Since(start)
							return
						}
				}
			case *tcell.EventMouse:
				if (ps.Buttons()&tcell.Button1 != 0) {
					// if the player clicked on a cell we will draw a cell there
					
					x, y := ps.Position()

					if s.NextCells[y-1][x-1] {
						s.NextCells[y-1][x-1] = false	
						s.S.SetContent( x, y, ' ', nil, tcell.Style{})
					}else {
						s.NextCells[y-1][x-1] = true
						s.S.SetContent( x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
					}


					s.S.Show()

				}
			}
		default: 
			// update status bar or something
			
		}
	}
}

func (s *Screen) Start() {
	
	defer s.S.Fini()
	
	s.S.EnableMouse()

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
					
				case tcell.KeyRune:
					if ev.Rune() == ' ' {
						// enter pause mode and stop the game
						s.PauseGame(eventChan)
					}
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

		timeElapsed := time.Since(s.StartTime) - s.PausedTime

		barContent := fmt.Sprintf("cell 1 size(%v, %v) FPS: %.2f Time: %.2f Alive: %v      ", s.Width, s.Height, float64(s.FrameCount)/ timeElapsed.Seconds(), timeElapsed.Seconds(), s.CellsAlive)

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

	s.CellsAlive = 0
	
	for i:= range s.Cells {
		for j:= range s.Cells[i] {
			// here do the trick
			s.UpdateCellState( i, j)

			if s.NextCells[i][j] {
				s.CellsAlive++
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

