import { imgPosterUrl, ky, PAGE, theme } from "../../utils.js"
import { computed, onMounted, ref } from "../../vue.js"

export const UpcomingTile = {
  props: {
    item: Object,
  },
  components: {},
  setup: (props) => {
    onMounted(() => {})

    const timeLeft = computed(() => {
      if (!props.item.Episode.AirDate) return null
      const date1 = new Date()
      const date2 = new Date(props.item.Episode.AirDate)
      const diffTime = Math.abs(date2 - date1)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      return diffDays
    })
    console.log(props.item)

    return {
      timeLeft,
      ky,
      imgPosterUrl,
    }
  },
  template: /* HTML */ `
    <div class="card">
      <div class="d-flex flex-row ">
        <img :src="imgPosterUrl(item.ImgPoster)" class="img-fluid rounded-start" alt="..." />
        <div class="card-body">
          <span> In {{ timeLeft }} days </span>
          <h5 class="card-title">
            <router-link :to="'/series/' + item.Episode.SeriesMId">
              {{ item.SeriesName }} <span class="text-muted">({{ item.Year }})</span>
            </router-link>
          </h5>
          <p class="card-text">{{ ky(item.Episode.Season, item.Episode.Episode) }}</p>
          <p class="card-text text-start">{{ (new Date(item.Episode.AirDate)).toDateString() }}</p>
          <p class="card-text">{{ item.Episode.Overview }}</p>
        </div>
      </div>
    </div>
  `,
}
