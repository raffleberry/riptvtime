import { notifyError } from "./components/Error.js";

const ENDPOINT = Object.freeze({
    FEED: () => { return '/api/series/feed' },
    SERIES_ALL: () => { return `/api/series` },
    SEARCH_SERIES: () => { return '/api/series/search' },
    SERIES_STATUS: (Id) => { return `/api/series/${Id}/status` },
    SERIES_ADD: () => { return `/api/series` },
    SERIES_REM: (Id) => { return `/api/series/${Id}` },
    SERIES_GET: (Id) => { return `/api/series/${Id}?full=1` },
})

export const apiSetStatus = async (Id, newStatus) => {
    let url = `${ENDPOINT.SERIES_STATUS(Id)}`
    try {
        const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ Status: newStatus })
        });

        if (!response.ok) {
            throw new Error(`Error, bad response from server: ${response.status} - ${response.statusText}`)
        }

        return {
            data: "OK",
            err: null
        }
    } catch (error) {
        notifyError(`Error updating status for series ID ${Id}: ${error.message}`);
        console.error('error updating status:', error);
        return {
            data: null,
            err: error
        }
    }
}

export const apiAddSeries = async (Id) => {
    try {
        const res = await fetch(ENDPOINT.SERIES_ADD(), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                MId: Id,
            })
        })

        if (res.status !== 200) {
            const errTxt = await res.text()
            return {
                err: new Error(`Error, response from server: ${res.status} - ${errTxt}`)
            }
        }

        const data = await res.json()

        return {
            data: data
        }

    } catch (error) {
        console.error(error)
        return {
            err: error
        }
    }
}

export const apiRemSeries = async (Id) => {
    try {

        const res = await fetch(ENDPOINT.SERIES_REM(Id), {
            method: 'DELETE',
        })

        if (res.status !== 200 && res.status !== 204 && res.status !== 404) {
            const errTxt = await res.text()
            return {
                err: new Error(`Error, response from server: ${res.status} - ${errTxt}`)
            }
        }
    } catch (error) {
        console.error(error)
        return {
            err: error
        }
    }

    return {}

}

export const apiSeriesTracked = async (Id) => { 
    try {

        const res = await fetch(ENDPOINT.SERIES_ALL())

        if (res.status !== 200) {
            const errTxt = await res.text()
            return {
                err: new Error(`Error, response from server: ${res.status} - ${errTxt}`)
            }
        }

        const data = await res.json()

        return {
            data: data,
        }

    } catch (error) {
        console.error(error)
        return {
            err: error
        }
    }

    return {}

}

