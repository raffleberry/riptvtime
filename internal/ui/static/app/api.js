import { notifyError } from "./components/Error.js";

const ENDPOINT = Object.freeze({
    FEED: () => { return '/api/series/feed' },
    SERIES_ALL: () => { return `/api/series` },
    SEARCH_SERIES: () => { return '/api/series/search' },
    SERIES_STATUS: (Id) => { return `/api/series/${Id}/status` },
    SERIES_ADD: () => { return `/api/series` },
    SERIES_REM: (Id) => { return `/api/series/${Id}` },
    SERIES_GET: (Id) => { return `/api/series/${Id}?full=1` },
    SERIES_EP_MARK: (Id) => { return `/api/series/episode` },
    SERIES_EP_UPNEXT: (Id) => { return `/api/series/${Id}/upnext` },
})

export const apiSearchTv = async (search, page) => {
    let url = `${ENDPOINT.SEARCH_SERIES()}?q=${search}&p=${page}`
    try {
        const res = await fetch(url);

        if (!res.ok) {
            throw new Error(`Error, response from server: ${res.status} - ${res.statusText}`)
        }

        const result = await res.json();
        return {
            data: result,
        }
    } catch (error) {
        console.error('Error fetching music data:', error);
        return {
            err: error
        }
    }
}

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

    } catch (error) {
        return {
            err: error
        }
    }
    return {}
}
export const apiGetStatus = async (mId) => {
    let url = `${ENDPOINT.SERIES_STATUS(mId)}`
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ Status: newStatus })
        });

        if (!response.ok) {
            throw new Error(`Error, bad response from server: ${response.status} - ${response.statusText}`)
        }

    } catch (error) {
        return {
            err: error
        }
    }
    return {}
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

        if (!res.ok) {
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

        if (!res.ok && res.status !== 404) {
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

        if (!res.ok) {
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

export const apiEpWatch = async (mId, sNo, epNo) => {
    try {
        const res = await fetch(ENDPOINT.SERIES_EP_MARK(), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                SeriesMId: mId,
                SeasonNo: sNo,
                EpisodeNo: epNo,
            })
        })

        if (!res.ok) {
            const errTxt = await res.text()
            return new Error(`Error, response from server: ${res.status} - ${errTxt}`)
        }
    } catch (error) {
        console.error(error)
        return error
    }
    return null
}

export const apiEpUnWatch = async (mId, sNo, epNo) => {
    try {
        const res = await fetch(ENDPOINT.SERIES_EP_MARK(), {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                SeriesMId: mId,
                SeasonNo: sNo,
                EpisodeNo: epNo,
            })
        })

        if (!res.ok) {
            const errTxt = await res.text()
            return new Error(`Error, response from server: ${res.status} - ${errTxt}`)

        }
    } catch (error) {
        console.error(error)
        return error

    }
    return null
}

export const apiEpUpNext = async (mId) => {
    try {
        const res = await fetch(ENDPOINT.SERIES_EP_UPNEXT(mId), {
            method: 'GET',
        })

        if (res.status !== 200) {

            if (res.status === 204) {
                return {
                    data: {
                        S: -1,
                        E: -1,
                    }
                }
            }

            const errTxt = await res.text()
            return { err: new Error(`Error, response from server: ${res.status} - ${errTxt}`) }
        }

        const data = await res.json()

        return {
            data: data
        }

    } catch (error) {
        console.error(error)
        return { err: error }

    }
    return {}
}


export const apiGetSeriesDetails = async (id) => {
    try {
        const res = await fetch(ENDPOINT.SERIES_GET(id));
        if (!res.ok) {
            return {
                err: new Error(`${res.status} - ${await res.text()}`)
            }
        }
        const data = await res.json();

        return {
            data: data
        }
    } catch (error) {
        console.error(error);
        return {
            err: new Error(`Error fetching series data:, ${error}`),
        }
    }
}