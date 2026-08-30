import { apiGetSeriesStats, apiGetSeriesStatsMyShows } from "../api.js"
import { notify } from "../components/Notify/Notify.js"
import { imgPosterUrl, imgPosterUrlFromMId } from "../utils.js"
import { onMounted, ref, watch } from "../vue.js"

const Stats = {
  props: {},
  components: {},
  setup: (props) => {
    const stats = ref({})
    const shows = ref([])
    const showAll = ref(false)

    const toggleShowAll = () => {
      showAll.value = !showAll.value
    }

    const fetchStats = async () => {
      const { data, err } = await apiGetSeriesStats()
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        return
      }
      stats.value = data
    }

    const showsLoading = ref(false)

    const fetchShows = async (limit) => {
      showsLoading.value = true
      const { data, err } = await apiGetSeriesStatsMyShows(limit)
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        return
      }
      shows.value = data
      showsLoading.value = false
    }

    watch(
      showAll,
      (v) => {
        if (v) {
          fetchShows(-1)
        } else {
          fetchShows(8)
        }
      },
      { immediate: true },
    )

    onMounted(() => {
      fetchStats()
    })

    const imgUrl = (id, mId) => {
      if (id) {
        return imgPosterUrl(id)
      }
      return imgPosterUrlFromMId(mId)
    }

    return {
      stats,
      shows,
      imgUrl,
      showAll,
      toggleShowAll,
      showsLoading,
    }
  },
  template: /* HTML */ `
    <div class="d-flex flex-grow-1 flex-column align-items-center">
      <div class="d-flex align-items-center">
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalHours) }}</h5>
            <p class="card-text">Hours Watched</p>
          </div>
        </div>
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalEpisodes) }}</h5>
            <p class="card-text">Episodes Watched</p>
          </div>
        </div>
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalShows) }}</h5>
            <p class="card-text">Shows Tracked</p>
          </div>
        </div>
      </div>
      <div class="mt-2 d-flex w-100 justify-content-center">
        <span class="mt-2 me-2">My Shows </span>
        <button
          type="button"
          @click="toggleShowAll"
          :class="!showAll ? 'text-bg-primary' : 'text-bg-secondary' "
          class="btn btn-link"
        >
          {{ showAll ? "hide":"all" }}
        </button>
      </div>
      <div v-if="!showsLoading" class="d-flex mt-2 flex-wrap justify-content-center">
        <div class="mx-2 my-1" v-for="show in shows" class="col-2" style="width: 180px">
          <img
            :src="imgUrl(show.Image, show.MId)"
            class="img-fluid object-fit-cover rounded-start"
            :alt="show.Name"
            loading="lazy"
          />
          <span>
            <router-link :to="'/series/' + show.MId"> {{show.Name}} </router-link>
            <span class="text-muted"> ({{show.Year}})</span>
          </span>
        </div>
      </div>
      <div v-else class="d-flex mt-5 justify-content-center">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
    </div>
  `,
}
export { Stats }
