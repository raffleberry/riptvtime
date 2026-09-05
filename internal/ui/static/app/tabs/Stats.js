import { apiGetSeriesFavs, apiGetSeriesStats, apiGetSeriesStatsMyShows } from "../api.js"
import { MsgType, notify } from "../components/Notify/Notify.js"
import { imgPosterUrl, imgPosterUrlFromMId } from "../utils.js"
import { onMounted, ref, watch } from "../vue.js"

const Stats = {
  props: {},
  components: {},
  setup: (props) => {
    const stats = ref({})
    const fetchStats = async () => {
      const { data, err } = await apiGetSeriesStats()
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        return
      }
      stats.value = data
    }

    const shows = ref([])
    const showsAll = ref(false)
    const showsLoading = ref(false)

    const ShowsAllTogl = () => {
      showsAll.value = !showsAll.value
    }

    const fetchShows = async (limit) => {
      showsLoading.value = true
      const { data, err } = await apiGetSeriesStatsMyShows(limit)
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        showsLoading.value = false
        return
      }
      shows.value = data
      showsLoading.value = false
    }

    watch(
      showsAll,
      (v) => {
        if (v) {
          fetchShows(-1)
        } else {
          fetchShows(8)
        }
      },
      { immediate: true },
    )

    const imgUrl = (id, mId) => {
      if (id) {
        return imgPosterUrl(id)
      }
      return imgPosterUrlFromMId(mId)
    }

    const favs = ref([])
    const favsLoading = ref(false)
    const favsAll = ref(false)

    const favsAllTogl = () => {
      favsAll.value = !favsAll.value
    }

    const fetchFavs = async (limit) => {
      favsLoading.value = true
      const { data, err } = await apiGetSeriesFavs(limit)
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        favsLoading.value = false
        return
      }
      favs.value = data
      favsLoading.value = false
    }

    watch(
      favsAll,
      (v) => {
        if (v) {
          fetchFavs(-1)
        } else {
          fetchFavs(8)
        }
      },
      { immediate: true },
    )

    onMounted(() => {
      fetchStats()
    })

    return {
      stats,

      shows,
      showsLoading,
      showsAll,
      ShowsAllTogl,
      imgUrl,

      favs,
      favsLoading,
      favsAll,
      favsAllTogl,
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
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.FavShows) }}</h5>
            <p class="card-text">Favourite Shows</p>
          </div>
        </div>
      </div>
      <div class="mt-2 d-flex w-100 justify-content-center">
        <span class="mt-2 me-2">My Shows </span>
        <button
          type="button"
          @click="ShowsAllTogl"
          :class="!showsAll ? 'text-bg-primary' : 'text-bg-secondary' "
          class="btn btn-link"
        >
          {{ showsAll ? "hide":"all" }}
        </button>
      </div>
      <div
        v-if="!showsLoading"
        class="d-flex mt-2 flex-wrap justify-content-center"
        style="min-height: 330px;"
      >
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
      <div
        v-else
        class="row w-100 mt-2 justify-content-center align-items-center"
        style="min-height: 330px;"
      >
        <div class="placeholder-glow h-100" role="status">
          <div class="placeholder col-12 h-100 bg-secondary rounded"></div>
        </div>
      </div>

      <div class="mt-2 d-flex w-100 justify-content-center">
        <span class="mt-2 me-2">Favourites</span>
        <button
          type="button"
          @click="favsAllTogl"
          :class="!favsAll ? 'text-bg-primary' : 'text-bg-secondary' "
          class="btn btn-link"
        >
          {{ favsAll ? "hide":"all" }}
        </button>
      </div>
      <div
        v-if="!favsLoading"
        class="d-flex mt-2 flex-wrap justify-content-center "
        style="min-height: 330px;"
      >
        <div class="mx-2 my-1" v-for="fav in favs" class="col-2" style="width: 180px">
          <img
            :src="imgUrl(fav.ImgPoster, fav.MId)"
            class="img-fluid object-fit-cover rounded-start"
            :alt="fav.Name"
            loading="lazy"
          />
          <span>
            <router-link :to="'/series/' + fav.MId"> {{fav.Name}} </router-link>
            <span class="text-muted"> ({{fav.Year}})</span>
          </span>
        </div>
      </div>
      <div
        v-else
        class="row w-100 mt-2 justify-content-center align-items-center"
        style="min-height: 330px;"
      >
        <div class="placeholder-glow h-100" role="status">
          <div class="placeholder col-12 h-100 bg-secondary rounded"></div>
        </div>
      </div>
    </div>
  `,
}
export { Stats }
