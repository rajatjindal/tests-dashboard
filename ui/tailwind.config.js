const colors = require('tailwindcss/colors');
const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./components/**/*.{js,vue,ts}",
		"./layouts/**/*.vue",
		"./pages/**/*.vue",
		"./plugins/**/*.{js,ts}",
		"./nuxt.config.{js,ts}",
		"./app.vue",
	],
	variants: {
		placeholderOpacity: ['responsive', 'focus', 'hover', 'active'],
	},
	safelist: [
		'col-span-1',
		'col-span-2',
		'col-span-3',
		'col-span-4',
		'col-span-5',
		'col-span-6',
		'col-span-7',
		'col-span-8',
		'col-span-9',
		'col-span-10',
		'col-span-11',
		'col-span-12',
		'grid-cols-1',
		'grid-cols-2',
		'grid-cols-3',
		'grid-cols-4',
		'grid-cols-5',
		'grid-cols-6',
		'grid-cols-7',
		'grid-cols-8',
		'grid-cols-9',
		'grid-cols-10',
		'grid-cols-11',
		'grid-cols-12',
		'grid-cols-13',
		'grid-cols-14',
		'grid-cols-15',
		'grid-cols-16',
		'grid-cols-17',
		'grid-cols-18',
		'grid-cols-19',
		'grid-cols-20'

	],
	theme: {
		colors: {
			// basic colors
			inherit: 'inherit',
			current: 'currentColor',
			transparent: 'transparent',
			black: '#000',
			white: '#fff',

			// core brand colors
			seagreen: '#34e8bd',
			oxfordblue: '#0d203f',
			rust: '#ef946c',
			lavender: '#bea7e5',
			colablue: '#0e8fdd',
			darkspace: '#213762',
			darkocean: '#0a455a',
			darkolive: '#1f7a8c',
			darkplum: '#525776',
			midgreen: '#1fbca0',
			midblue: '#345995',
			midgrey: '#e8eef6',
			lightgrey: '#d9dbe8',
			lightplum: '#d3c3d9',
			lightlavender: '#ece5ee',
			lightlemon: '#f9f7ee',


			// color themes as defined in Figma
			// https://www.figma.com/file/1ZoPrunzykFacq3ZJE5dIk/Fermyon-Playground?type=design&node-id=7029-178&mode=design&t=Y8QoAGkQ1FY7mgDI-4
			//
			// NOTE(bacongobbler): the naming scheme for this color system seems
			// a bit wonky. We should probably change it from a
			// "lightmode/darkmode system" (`bg-lightmode-blue-contrast2
			// dark:bg-darkmode-blue-contrast2`) to a numeric scaling system
			// (`bg-blue-100 dark:bg-blue-800`).
			//
			// https://tailwindcss.com/docs/customizing-colors
			lightmode: {
				primary: {
					DEFAULT: '#34e8bd',
					contrast1: '#063429',
					contrast2: '#1fbca0',
					contrast3: '#B8F1EB',
				},
				background: {
					DEFAULT: '#E8EEF6',
					contrast1: '#EEF4FC',
					contrast2: '#FAFAFC',
				},
				blue: {
					// NOTE(bacongobbler): the design document lists this as
					// "contrast0", but I'm defining it as the "default" color
					// for consistency.
					DEFAULT: '#14213D',
					contrast0: '#14213D',
					contrast1: '#8595BA',
					contrast2: '#99A6C3',
					contrast3: '#F1F5FE',
					contrast4: '#B3BDD7',
					contrast5: '#D6DEED',
					contrast6: '#707A8F',
					contrast7: '#4B5F8F',
				},
				titleunfocused: '#475470',
				red: {
					DEFAULT: '#CE5050',
					light: '#FED8D8',
					dark: '#773A48',
				},
				// 80% opacity
				backgroundblur: '#B3BDD7CC',
			},

			darkmode: {
				primary: {
					DEFAULT: '#34e8bd',
					contrast1: '#063429',
					contrast2: '#1fbca0',
					contrast3: '#177272',
				},
				background: {
					DEFAULT: '#0D203F',
					contrast1: '#152A4D',
					contrast2: '#061632',
				},
				blue: {
					// NOTE(bacongobbler): the design document lists this as
					// "contrast0", but I'm defining it as the "default" color
					// for consistency.
					DEFAULT: '#DBE3EF',
					contrast0: '#DBE3EF',
					contrast1: '#93A4C4',
					contrast2: '#586B8A',
					contrast3: '#25375F',
					contrast4: '#465B8B',
					contrast5: '#2D416E',
					contrast6: '#6379AA',
					contrast7: '#6C80AE',
				},
				titleunfocused: '#7E8BA0',
				red: {
					DEFAULT: '#CE5050',
					light: '#FED8D8',
					dark: '#ef4444',
				},
				// 80% opacity
				backgroundblur: '#030D1ECC',
			},
		},

		// typography
		fontFamily: {
			sans: [
				'Sen',
				'Europa',
				'Avenir',
				...defaultTheme.fontFamily.sans,
			],
			mono: [
				'SourceCodePro',
				...defaultTheme.fontFamily.mono,
			],
		},
	},
	plugins: [],
}