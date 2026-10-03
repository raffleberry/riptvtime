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
})

export const apiFetchFeed = async () => {
  let url = `${ENDPOINT.FEED()}`
  try {
    const res = await fetch(url)
    if (!res.ok) {
      throw new Error(`Error, response from server: ${res.status} - ${res.statusText}`)
    }
    const result = await res.json()
    return {
      data: result,
    }
  } catch (error) {
    console.error("Error fetching feed data:", error)
    return {
      err: error,
    }
  }
}

export const apiSearchTv = async (search, page) => {
  let url = `${ENDPOINT.SEARCH_SERIES()}?q=${search}&p=${page}`
  try {
    const res = await fetch(url)

    if (!res.ok) {
      throw new Error(`Error, response from server: ${res.status} - ${res.statusText}`)
    }

    const result = await res.json()
    return {
      data: result,
    }
  } catch (error) {
    console.error("Error Search data:", error)
    return {
      err: error,
    }
  }
}

export const apiSetStatus = async (Id, newStatus) => {
  let url = `${ENDPOINT.SERIES_STATUS(Id)}`
  try {
    const response = await fetch(url, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ Status: newStatus }),
    })

    if (!response.ok) {
      throw new Error(
        `Error, bad response from server: ${response.status} - ${response.statusText}`,
      )
    }
  } catch (error) {
    return {
      err: error,
    }
  }
  return {}
}
export const apiGetStatus = async (mId) => {
  let url = `${ENDPOINT.SERIES_STATUS(mId)}`
  try {
    const response = await fetch(url, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ Status: newStatus }),
    })

    if (!response.ok) {
      throw new Error(
        `Error, bad response from server: ${response.status} - ${response.statusText}`,
      )
    }
  } catch (error) {
    return {
      err: error,
    }
  }
  return {}
}

export const apiAddSeries = async (Id) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_ADD(), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        MId: Id,
      }),
    })

    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }

    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: error,
    }
  }
}

export const apiRemSeries = async (Id) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_REM(Id), {
      method: "DELETE",
    })

    if (!res.ok && res.status !== 404) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }
  } catch (error) {
    console.error(error)
    return {
      err: error,
    }
  }

  return {}
}

export const apiGetSeriesFavs = async (limit) => {
  try {
    if (!limit) limit = -1
    const res = await fetch(ENDPOINT.SERIES_FAVS(limit))
    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }
    const data = await res.json()
    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return error
  }

  return {}
}

export const apiSeriesTracked = async (Id) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_ALL())

    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }

    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: error,
    }
  }

  return {}
}

export const apiEpWatch = async (mId, eps) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_EP_MARK(), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        SeriesMId: mId,
        Episodes: eps,
      }),
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
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        SeriesMId: mId,
        SeasonNo: sNo,
        EpisodeNo: epNo,
      }),
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
      method: "GET",
    })

    if (res.status !== 200) {
      if (res.status === 204) {
        return {
          data: {
            S: -1,
            E: -1,
          },
        }
      }

      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }

    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return { err: error }
  }
  return {}
}

export const apiGetSeriesDetails = async (id) => {
  try {
    const res = await fetch(ENDPOINT.SERIES_GET(id))
    if (!res.ok) {
      return {
        err: new Error(`${res.status} - ${await res.text()}`),
      }
    }
    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: new Error(`Error fetching series data:, ${error}`),
    }
  }
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
        resolve(JSON.parse(xhr.responseText))
      } else {
        reject(xhr.statusText)
      }
    }

    xhr.onerror = () => {
      reject("Something went wrong")
    }
    xhr.open("POST", ENDPOINT.IMPORT_UPLOAD())
    xhr.send(formData)
  })
}

export const apiGetUnrImportData = async () => {
  try {
    const res = await fetch(ENDPOINT.IMPORT_UNRESOLVED())
    if (!res.ok) {
      return {
        err: new Error(`${res.status} - ${await res.text()}`),
      }
    }
    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: new Error(`Error fetching unmatched data:, ${error}`),
    }
  }
}

export const apiGetState = async () => {
  try {
    const res = await fetch(ENDPOINT.STATE())
    if (!res.ok) {
      return {
        err: new Error(`${res.status} - ${await res.text()}`),
      }
    }
    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: error,
    }
  }
}

export const apiResetState = async () => {
  try {
    const res = await fetch(ENDPOINT.STATE(), { method: "DELETE" })
    if (!res.ok) {
      return {
        err: new Error(`${res.status} - ${await res.text()}`),
      }
    }
    const data = await res.json()

    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return {
      err: error,
    }
  }
}
export const apiImportMatch = async (TvTimeSId, MId) => {
  try {
    const res = await fetch(ENDPOINT.IMPORT_MATCH(), {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        TvTimeSId: TvTimeSId,
        MId: MId,
      }),
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

export const apiIgnoreUnresolved = async (key) => {
  try {
    const res = await fetch(ENDPOINT.IMPORT_IGNORE(), {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ Key: key }),
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

export const apiGetSeriesStats = async () => {
  try {
    const res = await fetch(ENDPOINT.SERIES_STATS())
    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }
    const data = await res.json()
    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return error
  }
  return null
}

export const apiGetSeriesStatsMyShows = async (limit) => {
  try {
    if (!limit) limit = -1
    const res = await fetch(ENDPOINT.SERIES_STATS_MY_SHOWS(limit))
    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }
    const data = await res.json()
    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return error
  }
  return null
}

export const apiGetUpcoming = async () => {
  try {
    const res = await fetch(ENDPOINT.SERIES_UPCOMING())
    if (!res.ok) {
      const errTxt = await res.text()
      return {
        err: new Error(`Error, response from server: ${res.status} - ${errTxt}`),
      }
    }
    const data = await res.json()
    return {
      data: data,
    }
  } catch (error) {
    console.error(error)
    return error
  }
  return null
}
