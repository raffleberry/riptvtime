import { apiGetUnrImportData, apiIgnoreUnresolved } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { nextTick, onMounted, onUnmounted, ref } from "../../vue.js"
import { Match } from "./Match.js"

const Unresolved = {
  props: {},
  components: {
    Match,
  },
  setup: (props) => {
    const activeTab = ref("episodes")
    const data = ref({ Series: [], Episodes: [] })
    const loading = ref(false)
    const ignoring = ref(false)
    const selected = ref(null)
    const selectedIdx = ref(-1)

    const tabs = [
      { key: "episodes", label: "Episodes" },
      { key: "series", label: "Series" },
      { key: "movies", label: "Movies" },
    ]

    const cmpNum = (a, b) => (a ?? 0) - (b ?? 0)

    const sortUnresolved = (list) => {
      return [...(list || [])].sort(
        (a, b) =>
          cmpNum(a.TvTimeId, b.TvTimeId) ||
          cmpNum(a.Season, b.Season) ||
          cmpNum(a.Episode, b.Episode) ||
          String(a.Key || "").localeCompare(String(b.Key || "")),
      )
    }

    const fetchData = async () => {
      try {
        loading.value = true
        const { data: res, err } = await apiGetUnrImportData()
        if (err) {
          throw err
        }
        const fetched = res || {}
        data.value = {
          ...fetched,
          Series: sortUnresolved(fetched.Series),
          Episodes: sortUnresolved(fetched.Episodes),
        }
      } catch (error) {
        console.error("Error fetching unresolved data:", error)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchData()
    })

    const currentList = () => {
      if (activeTab.value === "episodes") return data.value.Episodes || []
      if (activeTab.value === "series") return data.value.Series || []
      return []
    }

    const onSelect = (item, idx) => {
      selected.value = item
      selectedIdx.value = idx
    }

    const onTabChange = (key) => {
      activeTab.value = key
      selected.value = null
      selectedIdx.value = -1
    }

    const onIgnore = async (item) => {
      if (!item?.Key) return
      if (activeTab.value === "movies") {
        // Movie ignore is not implemented on the backend yet (stub).
        notify(MsgType.Error, "Unresolved", new Error("Movie ignore is not supported yet"))
        return
      }
      try {
        ignoring.value = true
        const removedIdx = selectedIdx.value
        const err = await apiIgnoreUnresolved(item.Key)
        if (err) {
          throw err
        }
        if (activeTab.value === "episodes") {
          data.value.Episodes = data.value.Episodes.filter((e) => e.Key !== item.Key)
        } else if (activeTab.value === "series") {
          data.value.Series = data.value.Series.filter((s) => s.Key !== item.Key)
        }
        const list = currentList()
        if (list.length === 0) {
          selected.value = null
          selectedIdx.value = -1
        } else {
          // Select the next item, or the previous one if the removed item was last.
          const nextIdx = Math.min(Math.max(removedIdx, 0), list.length - 1)
          onSelect(list[nextIdx], nextIdx)
          scrollToSelected()
        }
      } catch (error) {
        console.error("Error ignoring unresolved item:", error)
        notify(MsgType.Error, "Unresolved", error)
      } finally {
        ignoring.value = false
      }
    }

    const onMatchDone = ({ TvTimeSId, MId }) => {
      // TODO: implement backend API call
      console.log("matchDone", TvTimeSId, MId)
      if (activeTab.value === "series") {
        data.value.Series = data.value.Series.filter((s) => s.TvTimeId !== TvTimeSId)
      }
      selected.value = null
      selectedIdx.value = -1
    }

    const rowRefs = ref([])

    const setRowRef = (el, idx) => {
      if (el) {
        rowRefs.value[idx] = el
      }
    }

    const scrollToSelected = async () => {
      await nextTick()
      const el = rowRefs.value[selectedIdx.value]
      if (el) {
        el.scrollIntoView({ block: "nearest", behavior: "smooth" })
      }
    }

    const onKeydown = (e) => {
      if (e.key === "ArrowDown") {
        e.preventDefault()
        const list = currentList()
        if (selectedIdx.value < list.length - 1) {
          onSelect(list[selectedIdx.value + 1], selectedIdx.value + 1)
          scrollToSelected()
        }
      } else if (e.key === "ArrowUp") {
        e.preventDefault()
        const list = currentList()
        if (selectedIdx.value > 0) {
          onSelect(list[selectedIdx.value - 1], selectedIdx.value - 1)
          scrollToSelected()
        }
      }
    }

    onMounted(() => {
      window.addEventListener("keydown", onKeydown)
    })

    onUnmounted(() => {
      window.removeEventListener("keydown", onKeydown)
    })

    return {
      activeTab,
      tabs,
      data,
      loading,
      ignoring,
      selected,
      selectedIdx,
      currentList,
      onSelect,
      onTabChange,
      onIgnore,
      onMatchDone,
      setRowRef,
    }
  },
  template: /* HTML */ `
    <div class="d-flex flex-column" style="height: 100%;">
      <ul class="nav nav-tabs mb-3">
        <li v-for="tab in tabs" :key="tab.key" class="nav-item">
          <button
            class="nav-link"
            :class="{ active: activeTab === tab.key }"
            @click="onTabChange(tab.key)"
          >
            {{ tab.label }}
          </button>
        </li>
      </ul>

      <div
        v-if="loading"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

      <div v-else class="d-flex flex-row flex-grow-1 gap-3" style="min-height: 0;">
        <div class="flex-grow-1 overflow-auto" style="min-width: 0;">
          <div v-if="currentList().length === 0" class="text-center text-muted py-4">
            No unresolved {{ activeTab }}
          </div>

          <table v-else class="table table-hover table-sm">
            <thead>
              <tr>
                <th>TvTimeId</th>
                <th>Name</th>
                <th v-if="activeTab === 'episodes'">Season</th>
                <th v-if="activeTab === 'episodes'">Episode</th>
                <th>Reason</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(item, idx) in currentList()"
                :key="item.Key"
                :ref="(el) => setRowRef(el, idx)"
                :class="{ 'table-primary': selectedIdx === idx }"
                style="cursor: pointer;"
                @click="onSelect(item, idx)"
              >
                <td>{{ item.TvTimeId }}</td>
                <td>{{ item.Name || item.SeriesName }}</td>
                <td v-if="activeTab === 'episodes'">{{ item.Season }}</td>
                <td v-if="activeTab === 'episodes'">{{ item.Episode }}</td>
                <td class="text-danger">{{ item.UnresolvedMsg }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div
          v-if="selected"
          class="flex-shrink-0 overflow-auto bg-body border-start"
          style="width: 400px; min-width: 300px;"
        >
          <div class="p-3">
            <div class="d-flex justify-content-between align-items-center mb-2">
              <h5 class="mb-0">Details</h5>
              <button
                class="btn btn-sm btn-outline-secondary"
                @click="selected = null; selectedIdx = -1"
              >
                <i class="bi bi-x"></i>
              </button>
            </div>

            <div v-if="activeTab === 'episodes'">
              <dl class="row mb-2">
                <dt class="col-sm-3">Key</dt>
                <dd class="col-sm-9">{{ selected.Key }}</dd>

                <dt class="col-sm-3">Series</dt>
                <dd class="col-sm-9">{{ selected.SeriesName }}</dd>

                <dt class="col-sm-3">TvTimeId</dt>
                <dd class="col-sm-9">{{ selected.TvTimeId }}</dd>

                <dt class="col-sm-3">EpisodeId</dt>
                <dd class="col-sm-9">{{ selected.TvTimeEId }}</dd>

                <dt class="col-sm-3">Season</dt>
                <dd class="col-sm-9">{{ selected.Season }}</dd>

                <dt class="col-sm-3">Episode</dt>
                <dd class="col-sm-9">{{ selected.Episode }}</dd>

                <dt class="col-sm-3">Reason</dt>
                <dd class="col-sm-9 text-danger">{{ selected.UnresolvedMsg }}</dd>
              </dl>
              <button
                class="btn btn-warning"
                :disabled="ignoring"
                @click="onIgnore(selected)"
              >
                Ignore
              </button>
            </div>

            <div v-if="activeTab === 'series'">
              <button
                class="btn btn-warning mb-2"
                :disabled="ignoring"
                @click="onIgnore(selected)"
              >
                Ignore
              </button>
              <Match :item="selected" @matchDone="onMatchDone"></Match>
            </div>

            <div v-if="activeTab === 'movies'">
              <p class="text-muted">Movie ignore is not supported yet.</p>
              <button class="btn btn-warning" disabled>Ignore</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
}

export { Unresolved }
