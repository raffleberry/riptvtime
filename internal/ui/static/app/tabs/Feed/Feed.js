import { ky } from "../../utils.js"
import { onMounted, storeToRefs } from "../../vue.js"
import { ConfirmOpts, openConfirm } from "./ConfirmOpts.js"
import { TvFeedTile } from "./TvFeedTile.js"
import { useFeedStore } from "./feedStore.js"

const Feed = {
  props: {},
  components: {
    TvFeedTile,
    ConfirmOpts,
  },
  setup: (props) => {
    const store = useFeedStore()
    const { feed, loading } = storeToRefs(store)
    const { fetchFeed } = store

    onMounted(() => {
      fetchFeed()
    })

    const { epMarkAndGetUpNext } = useFeedStore()

    const handleMarkUpNext = async (d) => {
      try {
        const resp = await openConfirm(`${d.Name} (${d.Year})`, `${ky(d.S, d.E)} mark as watched?`)
        if (resp) {
          await epMarkAndGetUpNext(d.MId, d.S, d.E)
        }
      } catch (err) {
        console.error(err)
        notify(MsgType.Error, "Feed", err)
      } finally {
      }
    }

    return {
      feed,
      loading,
      handleMarkUpNext,
    }
  },
  template: /* HTML */ `
    <ConfirmOpts></ConfirmOpts>
    <div class="container">
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
        <div v-for="tv in feed" :key="tv.ID" class="mb-3">
          <TvFeedTile :tv="tv" @markupnext="handleMarkUpNext"></TvFeedTile>
        </div>
      </div>
    </div>
  `,
}
export { Feed }
