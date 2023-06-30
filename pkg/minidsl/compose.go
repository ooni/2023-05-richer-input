package minidsl

//
// Stage composition
//

import "context"

// Compose composes two [Stage] together.
func Compose[A, B, C any](s1 Stage[A, B], s2 Stage[B, C]) Stage[A, C] {
	return &composeStage[A, B, C]{s1, s2}
}

type composeStage[A, B, C any] struct {
	s1 Stage[A, B]
	s2 Stage[B, C]
}

func (sx *composeStage[A, B, C]) Run(ctx context.Context, rtx Runtime, input Maybe[A]) Maybe[C] {
	result := sx.s1.Run(ctx, rtx, input)
	if result.Error != nil {
		return NewError[C](result.Error)
	}
	return sx.s2.Run(ctx, rtx, result)
}

// Compose3 composes 3 [Stage] together.
func Compose3[
	T0,
	T1,
	T2,
	T3 any,
](
	s1 Stage[T0, T1],
	s2 Stage[T1, T2],
	s3 Stage[T2, T3],
) Stage[T0, T3] {
	return Compose(s1, Compose(s2, s3))
}

// Compose4 composes 4 [Stage] together.
func Compose4[
	T0,
	T1,
	T2,
	T3,
	T4 any,
](
	s1 Stage[T0, T1],
	s2 Stage[T1, T2],
	s3 Stage[T2, T3],
	s4 Stage[T3, T4],
) Stage[T0, T4] {
	return Compose(s1, Compose3(s2, s3, s4))
}

// Compose5 composes 5 [Stage] together.
func Compose5[
	T0,
	T1,
	T2,
	T3,
	T4,
	T5 any,
](
	s1 Stage[T0, T1],
	s2 Stage[T1, T2],
	s3 Stage[T2, T3],
	s4 Stage[T3, T4],
	s5 Stage[T4, T5],
) Stage[T0, T5] {
	return Compose(s1, Compose4(s2, s3, s4, s5))
}