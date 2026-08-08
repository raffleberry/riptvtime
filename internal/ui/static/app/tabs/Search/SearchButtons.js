import { useSearchStore } from "./searchStore.js"
import { ref, storeToRefs, useRouter } from "../../vue.js"

const SearchButtons = {
  components: {},
  props: {
    term: String,
    pageCur: Number,
    pageTotal: Number,
  },
  setup(props) {
    const store = useSearchStore()
    const { getTotalPages } = store
    const { loading, resultsCnt } = storeToRefs(store)
    const router = useRouter()

    const page = ref(props.pageCur)

    const onPrvBtn = () => {
      if (page.value === 1) {
        return
      }
      page.value -= 1
      router.push({ path: "/search", query: { q: props.term, p: page.value } })
    }

    const onNxtBtn = () => {
      if (page.value >= getTotalPages()) {
        return
      }
      page.value += 1
      router.push({ path: "/search", query: { q: props.term, p: page.value } })
    }

    return {
      loading,
      resultsCnt,
      page,
      getTotalPages,

      onNxtBtn,
      onPrvBtn,
    }
  },

  template: /* HTML */ `
    <div class="d-flex flex-column justify-content-center align-items-center">
      <div class="input-group mt-3 justify-content-center">
        <button
          :disabled="page === 1 || loading"
          type="button"
          class="btn btn-outline-primary"
          @click="onPrvBtn"
        >
          <i class="bi bi-arrow-left"></i>
        </button>
        <div class="mx-3 d-flex flex-column align-items-center">
          <div>{{ resultsCnt }} Results</div>
          <div>{{ page }} / {{ getTotalPages() }}</div>
        </div>
        <button
          :disabled="page === getTotalPages() || loading"
          type="button"
          class="btn btn-outline-primary"
          @click="onNxtBtn"
        >
          <i class="bi bi-arrow-right"></i>
        </button>
      </div>
    </div>
  `,
}

export { SearchButtons }
