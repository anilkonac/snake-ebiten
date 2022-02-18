package main

// Rectangle compatible with float64 type parameters of the ebitenutil.DrawRect function.
type rectF64 struct {
	x, y          float64
	width, height float64
}

// Divide rectangle up to 4 ones according to where it is off-screen.
func (r rectF64) slice(rects *[]rectF64) {
	rightX := r.x + r.width
	bottomY := r.y + r.height

	if r.x < 0 { // left part is off-screen
		rectF64{r.x + screenWidth, r.y, -r.x, r.height}.slice(rects) // teleported left part
		rectF64{0, r.y, rightX, r.height}.slice(rects)               // part in the screen
		return
	} else if rightX > screenWidth { // right part is off-screen
		rectF64{0, r.y, rightX - screenWidth, r.height}.slice(rects) // teleported right part
		rectF64{r.x, r.y, screenWidth - r.x, r.height}.slice(rects)  // part in the screen
		return
	}

	if r.y < 0 { // upper part is off-screen
		rectF64{r.x, screenHeight + r.y, r.width, -r.y}.slice(rects) // teleported upper part
		rectF64{r.x, 0, r.width, bottomY}.slice(rects)               // part in the screen
		return
	} else if bottomY > screenHeight { // bottom part is off-screen
		rectF64{r.x, 0, r.width, bottomY - screenHeight}.slice(rects) // teleported bottom part
		rectF64{r.x, r.y, r.width, screenHeight - r.y}.slice(rects)   // part in the screen
		return
	}

	// Add sliced rectangle to the slice
	*rects = append(*rects, r)
}

func intersects(rectA, rectB rectF64) bool {
	aRightX := rectA.x + rectA.width
	bRightX := rectB.x + rectB.width
	aBottomY := rectA.y + rectA.height
	bBottomY := rectB.y + rectB.height

	if (rectA.x-rectB.x <= epsilon) && (aRightX-rectB.x <= epsilon) { // rectA is on the left side of rectB
		return false
	}

	if (rectA.x-bRightX >= -epsilon) && (aRightX-bRightX >= -epsilon) { // rectA is on the right side of rectB
		return false
	}

	if (rectA.y-rectB.y <= epsilon) && (aBottomY-rectB.y <= epsilon) { // rectA is above rectB
		return false
	}

	if (rectA.y-bBottomY >= -epsilon) && (aBottomY-bBottomY >= -epsilon) { // rectA is under rectB
		return false
	}

	return true
}
