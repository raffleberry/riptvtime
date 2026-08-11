import { MsgType, notify } from "../../components/Notify/Notify.js"
import { useTracked } from "../../stores/tracked.js"
import { imgBackdropUrl, imgPosterUrl, ky, TvStatus } from "../../utils.js"
import {
  computed,
  nextTick,
  onMounted,
  ref,
  storeToRefs,
  useRoute,
  useTemplateRef,
  watch,
} from "../../vue.js"
import { EpisodeOpts } from "./EpisodeOpts.js"
import { SeriesMarkPrev } from "./SeriesMarkPrev.js"
import { SeriesOpts } from "./SeriesOpts.js"
import { useSeriesStore } from "./seriesStore.js"

export const Series = {
  props: {},
  components: {
    SeriesOpts,
    EpisodeOpts,
    SeriesMarkPrev,
  },
  setup: (props) => {
    const r = useRoute()

    const seriesStore = useSeriesStore()
    const { loading, sd, SnWatchedEps, epWatchCnt } = storeToRefs(seriesStore)
    const Id = computed(() => r.params.id)

    const { epMarkWatched, fetchSeries, getWatchedEpsCnt } = seriesStore

    const { series } = storeToRefs(useTracked())

    const seriesId = ref(0)

    const status = computed(() => {
      let ob = series.value?.[Id.value]
      if (ob) {
        return ob.TrackingStatus
      }
      return TvStatus.NotWatching
    })

    const selectedEp = ref({})

    const eps = ref([])

    onMounted(() => {})
    watch(loading, async (val) => {
      if (!val) {
        if (r.hash) {
          await nextTick()
          let sel = document.querySelector(r.hash)
          if (sel) sel.scrollIntoView()
        }
      }
    })

    watch(sd, () => {
      console.log(sd.value)
      if (sd.value?.Name) {
        document.title = `${sd.value.Name} (${sd.value.Year}) - ${document.title}`
      }
    })

    const getStatusTxt = (s) => {
      switch (s) {
        case TvStatus.Watching:
          return "Watching"
        case TvStatus.Completed:
          return "Completed"
        case TvStatus.Stopped:
          return "Stopped"
        case TvStatus.UpToDate:
          return "Up To Date"
        default:
          return "Add"
      }
    }

    const statusCss = computed(() => {
      let card = "card"
      let icon = "bi"
      let btn = ""
      let pgbar = ""
      switch (status.value) {
        case TvStatus.Watching:
          card += " border border-warning"
          icon += " bi-bookmark-check"
          btn += " btn-warning"
          pgbar += "bg-warning"
          break
        case TvStatus.NotWatching:
          // cb += ' border border-primary'
          icon += " bi-bookmark-plus"
          btn += " btn-primary"
          pgbar += "bg-primary"
          break
        case TvStatus.Stopped:
          card += " border border-danger"
          icon += " bi-bookmark-x"
          btn += " btn-danger"
          pgbar += "bg-danger"
          break
        case TvStatus.Completed:
          card += " border border-success"
          icon += " bi-bookmark-check-fill"
          btn += " btn-success"
          pgbar += "bg-success"
          break
        case TvStatus.UpToDate:
          card += " border border-info"
          icon += " bi-bookmark-check-fill"
          btn += " btn-info"
          pgbar += "bg-info"
          break

        default:
          break
      }
      return {
        card,
        icon,
        btn,
        pgbar,
      }
    })

    const progress = computed(() => {
      if (sd.value?.EpisodesAired) {
        const deno = sd.value.EpisodesAired ?? 0
        if (deno === 0) {
          return 0
        }
        const num = getWatchedEpsCnt()
        return (num / deno) * 100
      }

      return 0
    })

    watch(
      Id,
      (id) => {
        if (id) {
          fetchSeries(id)
        }
      },
      { immediate: true },
    )

    var epOptEl = null
    var epMpEl = null

    onMounted(() => {
      const el = document.getElementById("episodeOpts")
      epOptEl = bootstrap.Offcanvas.getOrCreateInstance(el)
      const elmp = document.getElementById("seriesMarkPrev")
      epMpEl = bootstrap.Offcanvas.getOrCreateInstance(elmp)
      // el.addEventListener('show.bs.offcanvas', () => {
      // })

      // el.addEventListener('hidden.bs.offcanvas', () => {
      // })
    })

    const cnt = (s, e) => {
      return epWatchCnt.value[ky(s, e)] ?? 0
    }

    const cardStyle = computed(() => {
      return {
        height: "360px",
        backgroundImage: `url(${imgBackdropUrl(sd.value.ImgBackdrop)})`,
        backgroundSize: "cover",
        backgroundRepeat: "no-repeat",
        backgroundPosition: "center",
      }
    })

    const getPopEpsCnt = (s, e) => {
      let episodes = []
      let foundWatched = false
      let watchedCnt = 0
      let totalEps = 0

      for (let sNo = s; sNo >= 1; sNo--) {
        for (let eNo = e; eNo >= 1; eNo--) {
          totalEps += 1
          if (cnt(sNo, eNo) > 0) {
            foundWatched = true
            watchedCnt += 1
          }
          if (!foundWatched) {
            episodes.push({ S: sNo, E: eNo })
          }
        }
      }
      if (episodes.length + watchedCnt === totalEps) {
        return episodes
      }
      return [{ S: s, E: e }]
    }

    const openEpOpts = async (ep) => {
      if (cnt(ep.SeasonNumber, ep.EpisodeNumber) === 0) {
        const epsPopCount = getPopEpsCnt(ep.SeasonNumber, ep.EpisodeNumber)
        if (epsPopCount.length > 1) {
          eps.value = epsPopCount
          epMpEl.show()
        } else {
          await epMarkWatched(Number(Id.value), epsPopCount)
        }
      } else {
        selectedEp.value = ep
        epOptEl.show()
      }
    }

    const isAired = (dStr) => {
      return new Date() >= new Date(dStr)
    }

    const onClickAccordian = (num) => {
      document.querySelector("#season" + num).scrollIntoView()
    }

    return {
      ky,
      loading,
      sd,
      SnWatchedEps,
      epWatchCnt,
      statusCss,
      status,
      getStatusTxt,
      r,
      openEpOpts,
      selectedEp,
      cnt,
      progress,
      isAired,
      eps,
      imgPosterUrl,
      cardStyle,
      onClickAccordian,
    }
  },
  template: /* HTML */ `
    <SeriesMarkPrev :mid="sd.Id" :eps="eps"></SeriesMarkPrev>
    <EpisodeOpts :mid="sd.Id" :ep="selectedEp"></EpisodeOpts>
    <SeriesOpts :mid="sd.Id" :name="sd.Name" :year="sd.Year"></SeriesOpts>
    <div class="container-fluid">
      <div
        v-if="loading"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div v-else>
        <div :class="statusCss.card" style="z-index: 5;">
          <div class="card-body d-flex flex-column justify-content-between p-0" :style="cardStyle">
            <div class="align-self-end">
              <button
                :class="statusCss.btn"
                data-bs-toggle="offcanvas"
                data-bs-target="#seriesOpts"
                type="button"
                class="btn btn-sm p-2 d-inline-flex align-items-center justify-content-center rounded-top-0 rounded-end-0"
              >
                <i class="bi me-2" :class="statusCss.icon"></i>
                {{ getStatusTxt(status) }}
              </button>
            </div>

            <div class="txt-bg-blur">
              <div class="d-flex justify-content-between">
                <h3 class="card-title">{{ sd.Name }} ({{ sd.Year }})</h3>
              </div>

              <span class="card-text fst-italic ">{{ sd.Tagline }}</span><br />
              <p class="card-text">{{ sd.Overview }}</p>
              <!-- -->
            </div>
          </div>
          <div>
            <div
              class="progress rounded-top-0"
              role="progressbar"
              aria-label="Tv show progress"
              :aria-valuenow="progress"
              aria-valuemin="0"
              aria-valuemax="100"
            >
              <div
                :class="statusCss.pgbar"
                class="progress-bar"
                :style="{width: progress + '%'}"
              ></div>
            </div>
          </div>
        </div>
        <div class="accordion" id="seasons">
          <div v-for="sn in sd.Seasons" class="accordion-item">
            <h2 class="accordion-header">
              <button
                class="accordion-button"
                type="button"
                data-bs-toggle="collapse"
                :data-bs-target="'#season' + sn.SeasonNumber"
                @click="onClickAccordian(sn.SeasonNumber)"
              >
                <div class="d-flex w-100 justify-content-between pe-5">
                  <span>{{ sn.Name }}</span>
                  <span>{{SnWatchedEps[sn.SeasonNumber]?.length}}/{{sn.EpisodeCount}}</span>
                </div>
              </button>
            </h2>
            <div
              :class="{ 'show': '#season' + sn.SeasonNumber === r.hash }"
              :id="'season' + sn.SeasonNumber"
              class="accordion-collapse collapse"
              data-bs-parent="#seasons"
            >
              <div class="accordion-body">
                <div
                  v-for="(ep, idx) in sn.Episodes"
                  :class="idx % 2 === 0 ? 'bg-body' : 'bg-body-secondary'"
                  class="d-flex flex-row justify-content-between align-items-center"
                >
                  <div>{{ ep.SeasonNumber }}x{{ ep.EpisodeNumber }} - {{ ep.Name }}</div>
                  <div>
                    <span v-if="cnt(ep.SeasonNumber, ep.EpisodeNumber) > 1"
                      >{{ cnt(ep.SeasonNumber, ep.EpisodeNumber) }}x</span
                    >
                    <span v-if="!isAired(ep.AirDate)"
                      >{{ new Date(ep.AirDate).toDateString() }}</span
                    >
                    <button
                      :disabled="!isAired(ep.AirDate)"
                      class="btn"
                      :class="!isAired(ep.AirDate) ? 'border-0' : ''"
                      @click="() => { openEpOpts(ep) }"
                    >
                      <i
                        class="bi"
                        :class="cnt(ep.SeasonNumber, ep.EpisodeNumber) > 0 ? 'bi-check-circle-fill text-success' : 'bi-check-circle'"
                      ></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
}
