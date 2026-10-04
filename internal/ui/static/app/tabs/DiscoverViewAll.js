import { computed, onMounted, storeToRefs, useRoute, useRouter, watch } from "../vue.js"
import { SearchTile } from "./Search/SearchTile.js"
import { SearchTileOpts } from "./Search/SearchTileOpts.js"
import {
  DISCOVER_LISTS,
  DISCOVER_SORTS,
  defaultSortFor,
  useDiscoverStore,
} from "./discoverStore.js"

const DiscoverViewAll = {
  components: {
    SearchTile,
    SearchTileOpts,
  },
  setup() {
    const store = useDiscoverStore()
    const { viewLoading, viewTotalPages, viewTotalResults, viewPages } = storeToRefs(store)
    const { fetchView } = store

    const route = useRoute()
    const router = useRouter()

    const list = computed(() => route.query.view || "popular")
    const page = computed(() => parseInt(route.query.p || "1", 10))
    const sort = computed(() => route.query.sort_by || defaultSortFor(list.value))

    const title = computed(() => {
      const found = DISCOVER_LISTS.find((l) => l.key === list.value)
      return found ? found.title : list.value
    })

    const items = computed(() => viewPages.value[page.value] || [])

    const load = () => {
      fetchView(list.value, sort.value, page.value)
    }

    watch(
      () => [route.query.view, route.query.p, route.query.sort_by],
      () => load(),
      { immediate: true },
    )

    onMounted(() => load())

    const onSortChange = (event) => {
      router.push({
        path: "/discover",
        query: { view: list.value, p: 1, sort_by: event.target.value },
      })
    }

    const goPage = (p) => {
      if (p < 1 || p > viewTotalPages.value) return
      router.push({
        path: "/discover",
        query: { view: list.value, p: p, sort_by: sort.value },
      })
    }

    const goBack = () => {
      router.push({ path: "/discover" })
    }

    return {
      list,
      page,
      sort,
      title,
      items,
      viewLoading,
      viewTotalPages,
      viewTotalResults,
      DISCOVER_SORTS,
      onSortChange,
      goPage,
      goBack,
    }
  },
  template: /* HTML */ `
    <SearchTileOpts></SearchTileOpts>
    <div class="d-flex justify-content-between align-items-center px-3 my-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="goBack">
        <i class="bi bi-arrow-left"></i> Discover
      </button>
      <h4 class="m-0">{{ title }}</h4>
      <select
        class="form-select form-select-sm"
        style="width: auto;"
        :value="sort"
        @change="onSortChange"
      >
        <option v-for="s in DISCOVER_SORTS" :key="s" :value="s">Sort: {{ s }}</option>
      </select>
    </div>
    <div class="d-flex px-3 flex-column overflow-auto">
      <div
        v-if="viewLoading && items.length === 0"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div
        v-else-if="items.length === 0"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <h2>Nothing</h2>
      </div>
      <div v-else>
        <SearchTile class="mb-3" v-for="tv in items" :key="tv.Id" :tv="tv"></SearchTile>
      </div>
    </div>
    <div
      class="d-flex flex-column justify-content-center align-items-center my-2"
      v-if="viewTotalResults > 0"
    >
      <div class="input-group mt-3 justify-content-center">
        <button
          :disabled="page <= 1 || viewLoading"
          type="button"
          class="btn btn-outline-primary"
          @click="goPage(page - 1)"
        >
          <i class="bi bi-arrow-left"></i>
        </button>
        <div class="mx-3 d-flex flex-column align-items-center">
          <div>{{ viewTotalResults }} Results</div>
          <div>{{ page }} / {{ viewTotalPages }}</div>
        </div>
        <button
          :disabled="page >= viewTotalPages || viewLoading"
          type="button"
          class="btn btn-outline-primary"
          @click="goPage(page + 1)"
        >
          <i class="bi bi-arrow-right"></i>
        </button>
      </div>
    </div>
  `,
}

export { DiscoverViewAll }
