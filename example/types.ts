/* THIS FILE IS GENERATED, DO NOT EDIT */

export enum ExplicitEnum {
	SomeEnumValA = "a",
	SomeEnumValB = "b",
}

export interface ISomeData {
	x: number;
	Y: number;
	Z: string;
	/**
	 * An explicitly typed enum we define somewhere else
	 */
	W: ExplicitEnum;
}

export interface IOuter {
	inner: IInner;
}

export interface IInner {
	x: number | null | undefined;
	y: number | null;
}
