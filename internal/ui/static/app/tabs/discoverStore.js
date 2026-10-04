import { apiDiscoverTv } from "../api.js"
import { MsgType, notify } from "../components/Notify/Notify.js"
import { defineStore, ref } from "../vue.js"

export const DISCOVER_LISTS = Object.freeze([
  { key: "recommended", title: "Recommended For You" },
  { key: "airing_today", title: "Airing Today" },
  { key: "on_the_air", title: "On The Air" },
  { key: "top_rated", title: "Top Rated" },
  { key: "popular", title: "Popular" },
])

export const DISCOVER_SORTS = Object.freeze([
  "popularity.desc",
  "popularity.asc",
  "first_air_date.desc",
  "first_air_date.asc",
  "vote_average.desc",
  "vote_average.asc",
  "vote_count.desc",
  "vote_count.asc",
  "name.asc",
  "name.desc",
  "original_name.asc",
  "original_name.desc",
])

export const defaultSortFor = (list) => {
  if (list === "top_rated") return "vote_average.desc"
  return "popularity.desc"
}

export const useDiscoverStore = defineStore("discover", () => {
  // rails[listKey] = { items, totalPages, totalResults, loading, err }
  const rails = ref({})

  // view-all page state (one list at a time)
  const viewLoading = ref(false)
  const viewList = ref("")
  const viewSort = ref("")
  const viewTotalPages = ref(1)
  const viewTotalResults = ref(0)
  const viewPages = ref({})

  const fetchRail = async (list) => {
    const cur = rails.value[list]
    if (cur && (cur.loading || cur.items)) return
    rails.value = {
      ...rails.value,
      [list]: { items: null, totalPages: 1, totalResults: 0, loading: true, err: null },
    }
    try {
      const { data, err } = await apiDiscoverTv(list, 1)
      if (err) throw err
      rails.value = {
        ...rails.value,
        [list]: {
          items: data.Results || [],
          totalPages: data.TotalPages,
          totalResults: data.TotalResults,
          loading: false,
          err: null,
        },
      }
    } catch (error) {
      console.error(error)
      rails.value = {
        ...rails.value,
        [list]: { items: [], totalPages: 1, totalResults: 0, loading: false, err: error },
      }
      notify(MsgType.Error, "Discover", error.message)
    }
  }

  const fetchRails = async () => {
    await Promise.all(DISCOVER_LISTS.map((l) => fetchRail(l.key)))
  }

  const fetchView = async (list, sort, page) => {
    if (!sort) sort = defaultSortFor(list)
    if (!page) page = 1
    if (viewLoading.value) return
    if (viewList.value === list && viewSort.value === sort && viewPages.value[page]) {
      return
    }
    try {
      viewLoading.value = true
      if (viewList.value !== list || viewSort.value !== sort) {
        viewList.value = list
        viewSort.value = sort
        viewPages.value = {}
      }
      const { data, err } = await apiDiscoverTv(list, page, sort)
      if (err) throw err
      viewTotalPages.value = data.TotalPages
      viewTotalResults.value = data.TotalResults
      viewPages.value = {
        ...viewPages.value,
        [page]: data.Results || [],
      }
    } catch (error) {
      console.error(error)
      notify(MsgType.Error, "Discover", error.message)
    } finally {
      viewLoading.value = false
    }
  }

  return {
    rails,
    fetchRail,
    fetchRails,
    viewLoading,
    viewList,
    viewSort,
    viewTotalPages,
    viewTotalResults,
    viewPages,
    fetchView,
  }
})
