import { ofetch } from 'ofetch'
import { MomentInput } from 'moment'
import moment from 'moment'

const runtimeConfig = useRuntimeConfig().public.baseURL
const baseURL = "https://testreporter.usingspin.com"

export const myfetch = ofetch.create({
    baseURL: baseURL,
    retry: 0,
    async onRequest({ request, options }) {
        console.log('[fetch request]', `[${options.method} ${request}]`)
    },
    onRequestError: function ({ request, options, error }) { },
    onResponse: function ({ request, response, options }) { },
    onResponseError: function ({ request, response, options, error }) {
        console.error('[fetch response error]', response.status, response.body, error);
    },
})

export const formatDate = function (input: Date | string | undefined, format: string = 'MMM DD, HH:mm'): string {
	return moment(input as MomentInput).format(format)
}

export const humanDuration = function (input: number): string {
	return moment.duration(input, 'seconds').asMinutes().toFixed(0).toString() + ' mins'
}

