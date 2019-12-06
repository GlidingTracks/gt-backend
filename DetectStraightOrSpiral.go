package main

import (
	"fmt"
	"github.com/atedja/go-vector"
	"math"
)

func main() {

	var VectorList []vector.Vector
/*
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.05, 1.78, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2, 2, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2.67, 2.08, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.49, 2.24, 0}))
*/

	VectorList = append(VectorList, vector.NewWithValues([]float64{4.52, 2.41, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{5.71, 2.53, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{6.69, 3, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.62, 3.66, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.4, 4.62, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{9.02, 5.87, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{9.1, 6.92, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.64, 7.75, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.21, 8.28, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.82, 8.73, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.49, 9.32, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.1, 9.73, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{6.36, 10.09, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{5.42, 10.18, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{4.59, 10.03, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.83, 9.49, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.01, 8.98, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2.57, 7.98, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.89, 7.29, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.42, 6.52, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.15, 5.82, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.13, 5.19, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.53, 4.64, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2, 4, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2.66, 3.32, 0}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.47, 2.64, 0}))



	VectorList = append(VectorList, vector.NewWithValues([]float64{4.52, 2.41, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{5.71, 2.53, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{6.69, 3, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.62, 3.66, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.4, 4.62, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{9.02, 5.87, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{9.1, 6.92, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.64, 7.75, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{8.21, 8.28, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.82, 8.73, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.49, 9.32, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{7.1, 9.73, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{6.36, 10.09, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{5.42, 10.18, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{4.59, 10.03, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.83, 9.49, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.01, 8.98, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2.57, 7.98, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.89, 7.29, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.42, 6.52, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.15, 5.82, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.13, 5.19, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.53, 4.64, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2, 4, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{2.66, 3.32, 3}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{3.47, 2.64, 3}))


	fmt.Println(len(VectorList))
	fmt.Println(VectorList[0])
	fmt.Println(VectorList[25])
/*
	for i, j := 0, len(VectorList)-1; i < j; i, j = i+1, j-1 {
		VectorList[i], VectorList[j] = VectorList[j], VectorList[i]
	}
*/
	fmt.Println(VectorList[0])
	fmt.Println(VectorList[25])

	detectSpiralTopAndBot(VectorList)
/*


	VectorList = append(VectorList, vector.NewWithValues([]float64{0.85, 0.61, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1, 2, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.36, 2.87, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.26, 4.02, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.2, 5.51, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.65, 6.43, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.45, 7.57, 1}))
	VectorList = append(VectorList, vector.NewWithValues([]float64{1.43, 9.04, 1}))

	println(marchingLineCylinder(VectorList))
*/
}

func detectSpiralTopAndBot(points []vector.Vector) (vector.Vector) {
	var Clock [24]vector.Vector
	BasisV := vector.Subtract(points[1], points[0])
	CompetingV := vector.New(3)
	FirstHalfV := vector.New(3)
	SecondHalfV := vector.New(3)
	Bottom := vector.New(3)
	Stopper := 0

	Cycle := 0
	for i:=0; i<24; {
		Stopper++
		fmt.Println(Stopper)
		CompetingV = vector.Subtract(points[2+Cycle], points[1+Cycle])
		if angleBetweenVectors(CompetingV, BasisV) >= 15 {
			if FirstHalfV.Magnitude() == 0{
				fmt.Println("Halleluja", i)
				FirstHalfV = BasisV
			}
			if angleBetweenVectors(BasisV, FirstHalfV) >= 145 && SecondHalfV.Magnitude() == 0 {
				fmt.Println("Done1 ")
				SecondHalfV = BasisV
			}
			Clock[i]=points[2+Cycle]
			BasisV = CompetingV
			fmt.Println("success: ", (i+1), " And: ", Clock[i])

			if angleBetweenVectors(BasisV, SecondHalfV) >= 145 {
				fmt.Println("Done2 ")
				fmt.Println(points[i])
				//TEST CODE

				//TEST CODE


				for x:=0; x< Stopper;x++{
					fmt.Println(Bottom)
					Bottom = vector.Add(Bottom, points[x])
				}

				//Bottom.Scale(float64(1.0 / float64(len(points))))
				Bottom.Scale(float64(1.0 / float64(Stopper)))
				fmt.Println("Bottom: ", Bottom)
				break
			}

			i++
		}
		if Cycle == len(points)-3{
			fmt.Println("Shit")
			break
		}
		fmt.Println("Works", i)
		Cycle++
	}

	return CompetingV
}

func angleBetweenVectors(a, b vector.Vector) float64{
	Dot, _ := vector.Dot(a, b)
	return (math.Acos(Dot/((a.Magnitude())*(b.Magnitude())))) * (180/math.Pi)
}

func marchingLineCylinder(points []vector.Vector) []int {
	var RangeList []int

	for i:=0;i<= (len(points)/25);i++{
		if len(points) >= 50{
			if rollingLineCylinder(points){
				RangeList = append(RangeList, 0)
				RangeList = append(RangeList, 24)
				points = points[25:]
			} else {
				points = points[25:]
			}




		} else {
			if rollingLineCylinder(points){
				RangeList = append(RangeList, 0)
				RangeList = append(RangeList, len(points))
				return RangeList
			}else {
				return RangeList
			}
		}
	}

	if len(points) >= 50{
	} else {
		if rollingLineCylinder(points){
			return append(RangeList, 0)
			return append(RangeList, len(points))
		}else {
			return RangeList
		}
	}

	return RangeList
}

func rollingLineCylinder(points []vector.Vector) bool {
	for i:=1; i<= (len(points)-1);i++{
		fmt.Println("hurray ")
		if !(lineCylinder(points[0], points[len(points)-1], points[i])) {
			fmt.Println("shit")
			return false
		}
	}
	fmt.Println("hey")
	return true
}

func lineCylinder(p1, p2, q vector.Vector) bool {
	var Range float64 = 1

	firstClear, _ := vector.Dot(vector.Subtract(q, p1), vector.Subtract(p2, p1))
	secondClear, _ := vector.Dot(vector.Subtract(q,p2), vector.Subtract(p2, p1))

	thirdClear, _ := vector.Cross(vector.Subtract(q, p1), vector.Subtract(p2, p1))
	fourthClear := vector.Subtract(p2, p1)
	fifthClear := thirdClear.Magnitude()/fourthClear.Magnitude()

	if  firstClear >= 0 && secondClear <= 0 && fifthClear <= Range {
		return true
	} else {
		return false
	}
}