package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	halfUnitLength = unitLength / 2.0
)

type unit struct {
	centerX   float64
	centerY   float64
	direction uint8
}

func (u *unit) moveUp(dist float64) {
	u.centerY -= dist

	// Teleport if center is off the screen
	if u.centerY < 0 {
		u.centerY += screenWidth
	}
}

func (u *unit) moveDown(dist float64) {
	u.centerY += dist

	// Teleport if center is off the screen
	if u.centerY > screenHeight {
		u.centerY -= screenWidth
	}
}

func (u *unit) moveLeft(dist float64) {
	u.centerX -= dist

	// Teleport if center is off the screen
	if u.centerX < 0 {
		u.centerX += screenWidth
	}
}

func (u *unit) moveRight(dist float64) {
	u.centerX += dist

	// Teleport if center is off the screen
	if u.centerX > screenWidth {
		u.centerX -= screenWidth
	}
}

// Checks if unit is off the screen
// Returns which axes should be sliced, and position of center on sliced axes
func (u unit) checkOffScreen() (slicedV, slicedH bool, locCenterX, locCenterY uint8) {
	if u.centerX-halfUnitLength < 0 {
		slicedV = true
		locCenterX = directionLeft
	} else if u.centerX+halfUnitLength > screenWidth {
		slicedV = true
		locCenterX = directionRight
	}

	if u.centerY-halfUnitLength < 0 {
		slicedH = true
		locCenterY = directionUp
	} else if u.centerY+halfUnitLength > screenHeight {
		slicedH = true
		locCenterY = directionDown
	}

	return
}

func (u unit) draw(screen *ebiten.Image, color color.Color) {
	slicedV, slicedH, locCenterX, locCenterY := u.checkOffScreen()
	if slicedV && slicedH { // unit is on one of the corner
		var rectangles [4]*screenRect
		switch {
		// Center is near top left corner
		case locCenterX == directionLeft && locCenterY == directionUp:
			rectangles = u.sliceTopLeft()

		// Center is near top right corner
		case locCenterX == directionRight && locCenterY == directionUp:
			rectangles = u.sliceTopRight()

		// Center is near bottom left corner
		case locCenterX == directionLeft && locCenterY == directionDown:
			rectangles = u.sliceBottomLeft()

		// Center is near bottom right corner
		case locCenterX == directionRight && locCenterY == directionDown:
			rectangles = u.sliceBottomRight()
		}

		// Draw rectangles
		for _, rect := range rectangles {
			if rect != nil {
				rect.draw(screen, color)
			} else {
				panic("Slicing to 4 rectangles is not successful!")
			}
		}
	} else if slicedV { // unit is on vertical edges
		var rectLeft, rectRight *screenRect
		switch locCenterX {
		case directionLeft:
			rectLeft, rectRight = u.sliceLeft()
		case directionRight:
			rectLeft, rectRight = u.sliceRight()
		default:
			panic("Wrong center X location!")
		}

		// Draw rectangles
		if rectLeft != nil && rectRight != nil {
			rectLeft.draw(screen, color)
			rectRight.draw(screen, color)
		} else {
			panic("Vertical slicing is not successful.")
		}
	} else if slicedH { // unit is on horizontal edges
		var rectUp, rectDown *screenRect
		switch locCenterY {
		case directionUp:
			rectUp, rectDown = u.sliceTop()
		case directionDown:
			rectUp, rectDown = u.sliceBottom()
		default:
			panic("Wrong center X location!")
		}

		// Draw rectangles
		if rectUp != nil && rectDown != nil {
			rectUp.draw(screen, color)
			rectDown.draw(screen, color)
		} else {
			panic("Horizontal slicing is not successful.")
		}
	} else { // unit is inside the screen
		newScreenRect(u.centerX-halfUnitLength, u.centerY-halfUnitLength, unitLength, unitLength).draw(screen, color)
	}
}

// Slice unit whose center is near top left to 4 rectangles
func (u unit) sliceTopLeft() (rects [4]*screenRect) {
	widthBigger := u.centerX + halfUnitLength
	widthSmaller := halfUnitLength - u.centerX
	heightBigger := u.centerY + halfUnitLength
	heightSmaller := halfUnitLength - u.centerY

	// Define rectangles
	rects[0] = newScreenRect(0, 0, widthBigger, heightBigger)
	rects[1] = newScreenRect(screenWidth-widthSmaller, 0, widthSmaller, heightBigger)
	rects[2] = newScreenRect(0, screenHeight-heightSmaller, widthBigger, heightSmaller)
	rects[3] = newScreenRect(screenWidth-widthSmaller, screenHeight-heightSmaller, widthSmaller, heightSmaller)

	return
}

// Slice unit whose center is near top right to 4 rectangles
func (u unit) sliceTopRight() (rects [4]*screenRect) {
	widthBigger := halfUnitLength + screenWidth - u.centerX
	widthSmaller := unitLength - widthBigger
	heightBigger := u.centerY + halfUnitLength
	heightSmaller := unitLength - heightBigger

	// Define rectangles
	rects[0] = newScreenRect(0, 0, widthSmaller, heightBigger)
	rects[1] = newScreenRect(screenWidth-widthBigger, 0, widthBigger, heightBigger)
	rects[2] = newScreenRect(0, screenHeight-heightSmaller, widthSmaller, heightSmaller)
	rects[3] = newScreenRect(screenWidth-widthBigger, screenHeight-heightSmaller, widthBigger, heightSmaller)

	return
}

// Slice unit whose center is near bottom left to 4 rectangles
func (u unit) sliceBottomLeft() (rects [4]*screenRect) {
	widthBigger := u.centerX + halfUnitLength
	widthSmaller := halfUnitLength - u.centerX
	heightBigger := halfUnitLength + screenHeight - u.centerY
	heightSmaller := unitLength - heightBigger

	// Define rectangles
	rects[0] = newScreenRect(0, 0, widthBigger, heightSmaller)
	rects[1] = newScreenRect(screenWidth-widthSmaller, 0, widthSmaller, heightSmaller)
	rects[2] = newScreenRect(0, screenHeight-heightBigger, widthBigger, heightBigger)
	rects[3] = newScreenRect(screenWidth-widthSmaller, screenHeight-heightBigger, widthSmaller, heightBigger)

	return
}

// Slice unit whose center is near bottom right to 4 rectangles
func (u unit) sliceBottomRight() (rects [4]*screenRect) {
	widthBigger := halfUnitLength + screenWidth - u.centerX
	widthSmaller := unitLength - widthBigger
	heightBigger := halfUnitLength + screenHeight - u.centerY
	heightSmaller := unitLength - heightBigger

	// Define rectangles
	rects[0] = newScreenRect(0, 0, widthSmaller, heightSmaller)
	rects[1] = newScreenRect(screenWidth-widthBigger, 0, widthBigger, heightSmaller)
	rects[2] = newScreenRect(0, screenHeight-heightBigger, widthSmaller, heightBigger)
	rects[3] = newScreenRect(screenWidth-widthBigger, screenHeight-heightBigger, widthBigger, heightBigger)

	return
}

// Divide unit to two rectangles if center is near left edge
func (u unit) sliceLeft() (rectLeft, rectRight *screenRect) {
	widthBigger := u.centerX + halfUnitLength
	widthSmaller := unitLength - widthBigger
	yLoc := u.centerY - halfUnitLength

	rectLeft = newScreenRect(0, yLoc, widthBigger, unitLength)
	rectRight = newScreenRect(screenWidth-widthSmaller, yLoc, widthSmaller, unitLength)
	return
}

// Divide unit to two rectangles if center is near right edge
func (u unit) sliceRight() (rectLeft, rectRight *screenRect) {
	widthBigger := halfUnitLength + screenWidth - u.centerX
	widthSmaller := unitLength - widthBigger
	yLoc := u.centerY - halfUnitLength

	rectLeft = newScreenRect(0, yLoc, widthSmaller, unitLength)
	rectRight = newScreenRect(screenWidth-widthBigger, yLoc, widthBigger, unitLength)
	return
}

// Divide unit to two rectangles if center is near top edge
func (u unit) sliceTop() (rectTop, rectBottom *screenRect) {
	heightBigger := u.centerY + halfUnitLength
	heightSmaller := unitLength - heightBigger
	xLoc := u.centerX - halfUnitLength

	rectTop = newScreenRect(xLoc, 0, unitLength, heightBigger)
	rectBottom = newScreenRect(xLoc, screenHeight-heightSmaller, unitLength, heightSmaller)
	return
}

// Divide unit to two rectangles if center is near bottom edge
func (u unit) sliceBottom() (rectTop, rectBottom *screenRect) {
	heightBigger := halfUnitLength + screenHeight - u.centerY
	heightSmaller := unitLength - heightBigger
	xLoc := u.centerX - halfUnitLength

	rectTop = newScreenRect(xLoc, 0, unitLength, heightSmaller)
	rectBottom = newScreenRect(xLoc, screenHeight-heightBigger, unitLength, heightBigger)
	return
}
