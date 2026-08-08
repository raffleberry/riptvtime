import { Navigation } from "../../components/Navigation.js"
import { onMounted, ref, RouterView, useRoute, useRouter, watch } from "../../vue.js"
import { useSearchStore } from "./searchStore.js"

const SearchBox = {
  components: {},
  props: {
    term: String,
  },

  setup(props) {
    const router = useRouter()
    const route = useRoute()

    const store = useSearchStore()
    const { onSearch } = store

    const searchTerm = ref(props.term || "")

    // watch(
    //   () => {
    //     return {
    //       q: route.query.q,
    //       p: route.query.p,
    //     }
    //   },
    //   (val) => {
    //     if (val.q && val.p) {
    //       searchTerm.value = val.q
    //       onSearch(val.q, val.p)
    //     } else if (val.q) {
    //       searchTerm.value = val.q
    //       onSearch(val.q, 1)
    //     }
    //   },
    //   { immediate: true },
    // )

    const handleSearch = () => {
      let st = searchTerm.value.trim()
      if (st !== "") {
        console.log(st)
        router.push({ path: "/search", query: { q: st, p: 1 } })
      }
    }

    return {
      searchTerm,
      handleSearch,
    }
  },

  template: /* HTML */ `
    <div class="row">
      <div class="input-group search-container">
        <input
          v-model="searchTerm"
          type="text"
          class="form-control search-input"
          placeholder="Search..."
          @keyup.enter="handleSearch"
        />
        <button
          class="btn btn-outline-primary"
          type="button"
          id="searchButton"
          @click="handleSearch"
        >
          <i class="bi bi-search"></i>
        </button>
      </div>
    </div>
  `,
}

export { SearchBox }
