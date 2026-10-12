import { useSearchStore } from "./searchStore.js"
import { computed, ref, storeToRefs, useRouter, watch } from "../../vue.js"

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

    const page = ref(Number(props.pageCur) || 1)
    watch(
      () => props.pageCur,
      (v) => {
        page.value = Number(v) || 1
      },
    )

    const totalPages = computed(() => props.pageTotal || getTotalPages() || 1)

    const onPrvBtn = () => {
      if (page.value <= 1) {
        return
      }
      router.push({ path: "/search", query: { q: props.term, p: page.value - 1 } })
    }

    const onNxtBtn = () => {
      if (page.value >= totalPages.value) {
        return
      }
      router.push({ path: "/search", query: { q: props.term, p: page.value + 1 } })
    }

    return {
      loading,
      resultsCnt,
      page,
      totalPages,
      getTotalPages,

      onNxtBtn,
      onPrvBtn,
    }
  },

  template: /* HTML */ `
    <div class="d-flex flex-column justify-content-center align-items-center">
      <div class="input-group mt-3 justify-content-center">
        <button
          :disabled="page <= 1 || loading"
          type="button"
          class="btn btn-outline-primary"
          @click="onPrvBtn"
        >
          <i class="bi bi-arrow-left"></i>
        </button>
        <div class="mx-3 d-flex flex-column align-items-center">
          <div>{{ resultsCnt }} Results</div>
          <div>{{ page }} / {{ totalPages }}</div>
        </div>
        <button
          :disabled="page >= totalPages || loading"
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
