import { base } from "./base.js"

const ENDPOINT = Object.freeze({
  FEED: () => {
    return `${base}/api/series/feed`
  },
  SERIES_ALL: () => {
    return `${base}/api/series`
  },
  SEARCH_SERIES: () => {
    return `${base}/api/series/search`
  },
  SERIES_STATUS: (Id) => {
    return `${base}/api/series/${Id}/status`
  },
  SERIES_ADD: () => {
    return `${base}/api/series`
  },
  SERIES_REM: (Id) => {
    return `${base}/api/series/${Id}`
  },
  SERIES_GET: (Id) => {
    return `${base}/api/series/${Id}?full=1`
  },
  SERIES_EP_MARK: (Id) => {
    return `${base}/api/series/episode`
  },
  SERIES_EP_UPNEXT: (Id) => {
    return `${base}/api/series/${Id}/upnext`
  },
  SERIES_UPCOMING: () => {
    return `${base}/api/series/upcoming`
  },
  SERIES_FAVS: (limit) => {
    return `${base}/api/series/favs?limit=${limit}`
  },
  IMPORT_UPLOAD: () => {
    return `${base}/api/import/upload`
  },
  IMPORT_UNRESOLVED: () => {
    return `${base}/api/import/unresolved`
  },
  IMPORT_MATCH: () => {
    return `${base}/api/import/match`
  },
  IMPORT_RESOLVE: () => {
    return `${base}/api/import/resolve`
  },
  IMPORT_IGNORE: () => {
    return `${base}/api/import/ignore`
  },
  STATE: () => {
    return `${base}/api/state`
  },
  SERIES_STATS: () => {
    return `${base}/api/series/stats`
  },
  SERIES_STATS_MY_SHOWS: (limit) => {
    return `${base}/api/series/stats/my?limit=${limit}`
  },
  DISCOVER_TV: (list, page, sortBy) => {
    let url = `${base}/api/discover/tv?list=${list}&p=${page}`
    if (sortBy) {
      url += `&sort_by=${encodeURIComponent(sortBy)}`
    }
    return url
  },
})

// Shared fetch helpers. Every api* function resolves to { data, err } so
// callers have a single contract to handle. Void endpoints resolve to {}.
const errOf = async (res) => {
  let txt = ""
  try {
    txt = await res.text()
  } catch {
    txt = res.statusText
  }
  return new Error(`Error, response from server: ${res.status} - ${txt}`)
}

const getJson = async (url, fetchOpts) => {
  try {
    const res = await fetch(url, fetchOpts)
    if (!res.ok) {
      return { err: await errOf(res) }
    }
    // 204 No Content has no body to parse.
    if (res.status === 204) {
      return { data: null }
    }
    const data = await res.json()
    return { data }
  } catch (error) {
    console.error("Error fetching:", url, error)
    return { err: error }
  }
}

const sendJson = async (url, method, body) => {
  try {
    const res = await fetch(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: body === undefined ? undefined : JSON.stringify(body),
    })
    if (!res.ok) {
      return { err: await errOf(res) }
    }
    const ct = res.headers.get("content-type") || ""
    if (ct.includes("application/json")) {
      return { data: await res.json() }
    }
    return {}
  } catch (error) {
    console.error("Error fetching:", url, error)
    return { err: error }
  }
}

export const apiFetchFeed = async () => {
  return getJson(ENDPOINT.FEED())
}

export const apiSearchTv = async (search, page) => {
  let url = `${ENDPOINT.SEARCH_SERIES()}?q=${encodeURIComponent(search)}&p=${page}`
  return getJson(url)
}

export const apiSetStatus = async (Id, newStatus) => {
  return sendJson(ENDPOINT.SERIES_STATUS(Id), "PUT", { Status: newStatus })
}

export const apiGetStatus = async (mId) => {
  return getJson(ENDPOINT.SERIES_STATUS(mId))
}

export const apiAddSeries = async (Id) => {
  return sendJson(ENDPOINT.SERIES_ADD(), "POST", { MId: Id })
}

export const apiRemSeries = async (Id) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_REM(Id), { method: "DELETE" })
    if (!res.ok && res.status !== 404) {
      return { err: await errOf(res) }
    }
    return {}
  } catch (error) {
    console.error(error)
    return { err: error }
  }
}

export const apiGetSeriesFavs = async (limit) => {
  if (!limit) limit = -1
  return getJson(ENDPOINT.SERIES_FAVS(limit))
}

export const apiSeriesTracked = async () => {
  return getJson(ENDPOINT.SERIES_ALL())
}

export const apiEpWatch = async (mId, eps) => {
  return sendJson(ENDPOINT.SERIES_EP_MARK(), "POST", {
    SeriesMId: mId,
    Episodes: eps,
  })
}

export const apiEpUnWatch = async (mId, sNo, epNo) => {
  return sendJson(ENDPOINT.SERIES_EP_MARK(), "PUT", {
    SeriesMId: mId,
    SeasonNo: sNo,
    EpisodeNo: epNo,
  })
}

export const apiEpUpNext = async (mId) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_EP_UPNEXT(mId), { method: "GET" })
    if (res.status === 204) {
      return { data: { S: -1, E: -1 } }
    }
    if (!res.ok) {
      return { err: await errOf(res) }
    }
    const data = await res.json()
    return { data }
  } catch (error) {
    console.error(error)
    return { err: error }
  }
}

export const apiGetSeriesDetails = async (id) => {
  return getJson(ENDPOINT.SERIES_GET(id))
}

export const apiUploadImportZip = async (formData, progress) => {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable) {
        const percentage = (event.loaded / event.total) * 100
        progress.value = percentage.toFixed(2)
      }
    }
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          resolve(JSON.parse(xhr.responseText))
        } catch {
          resolve({})
        }
      } else {
        reject(new Error(xhr.statusText || `Upload failed: ${xhr.status}`))
      }
    }

    xhr.onerror = () => {
      reject(new Error("Something went wrong"))
    }
    xhr.open("POST", ENDPOINT.IMPORT_UPLOAD())
    xhr.send(formData)
  })
}

export const apiGetUnrImportData = async () => {
  return getJson(ENDPOINT.IMPORT_UNRESOLVED())
}

export const apiGetState = async () => {
  return getJson(ENDPOINT.STATE())
}

export const apiResetState = async () => {
  try {
    const res = await fetch(ENDPOINT.STATE(), { method: "DELETE" })
    if (!res.ok) {
      return { err: await errOf(res) }
    }
    return { data: await res.json() }
  } catch (error) {
    console.error(error)
    return { err: error }
  }
}

export const apiImportMatch = async (TvTimeSId, MId) => {
  return sendJson(ENDPOINT.IMPORT_MATCH(), "PUT", { TvTimeSId, MId })
}

export const apiIgnoreUnresolved = async (key) => {
  return sendJson(ENDPOINT.IMPORT_IGNORE(), "PUT", { Key: key })
}

export const apiGetSeriesStats = async () => {
  return getJson(ENDPOINT.SERIES_STATS())
}

export const apiGetSeriesStatsMyShows = async (limit) => {
  if (!limit) limit = -1
  return getJson(ENDPOINT.SERIES_STATS_MY_SHOWS(limit))
}

export const apiGetUpcoming = async () => {
  return getJson(ENDPOINT.SERIES_UPCOMING())
}

export const apiDiscoverTv = async (list, page, sortBy) => {
  if (!page) page = 1
  return getJson(ENDPOINT.DISCOVER_TV(list, page, sortBy))
}
