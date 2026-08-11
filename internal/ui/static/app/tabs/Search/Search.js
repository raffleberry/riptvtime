import { SearchButtons } from "./SearchButtons.js"
import { PAGE } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, useRoute, watch } from "../../vue.js"
import { SearchBox } from "./SearchBox.js"
import { SearchTile } from "./SearchTile.js"
import { SearchTileOpts } from "./SearchTileOpts.js"
import { useSearchStore } from "./searchStore.js"

const Search = {
  props: {},
  components: {
    SearchTileOpts,
    SearchTile,
    SearchBox,
    SearchButtons,
  },
  setup: (props) => {
    const store = useSearchStore()

    const { loading, results, resultsCnt } = storeToRefs(store)
    const { onSearch, getTotalPages } = store

    const r = useRoute()
    const curPath = computed(() => r.path)

    const pageCur = computed(() => r.query.p || 1)

    const searchTerm = computed(() => r.query.q || "")

    watch(
      () => {
        return {
          q: r.query.q,
          p: r.query.p,
        }
      },
      (val) => {
        if (val.q) {
          document.title = `${val.q} - ${document.title}`
        }
        if (val.q && val.p) {
          onSearch(val.q, val.p)
        } else if (val.q) {
          onSearch(val.q, 1)
        }
      },
      { immediate: true },
    )

    onMounted(() => {})

    return {
      loading,
      searchTerm,
      results,
      resultsCnt,
      pageCur,
      curPath,
      getTotalPages,
      PAGE,
    }
  },
  template: /* HTML */ `
    <SearchTileOpts></SearchTileOpts>
    <SearchBox :term="searchTerm" class="my-2 ps-3" v-if="curPath === PAGE.SEARCH.path">
    </SearchBox>
    <div class="d-flex px-3 flex-column overflow-auto">
      <div
        v-if="loading"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div
        v-else-if="resultsCnt === 0"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <h2>Nothing</h2>
      </div>
      <div v-else>
        <SearchTile class="mb-3" v-for="tv in results[pageCur]" :key="tv.Id" :tv="tv"></SearchTile>
      </div>
    </div>
    <SearchButtons
      :term="searchTerm"
      :pageCur="pageCur"
      :pageTotal="getTotalPages()"
      class="my-2"
      v-if="curPath === PAGE.SEARCH.path && resultsCnt > 0"
    >
    </SearchButtons>
  `,
}
export { Search }
