import { nextTick, onBeforeUnmount, watch, type Ref } from 'vue'

type UseInfiniteScrollOptions = {
  target: Ref<HTMLElement | null>
  enabled: Ref<boolean>
  onLoadMore: () => void | Promise<void>
  rootMargin?: string
}

export function useInfiniteScroll(options: UseInfiniteScrollOptions) {
  let observer: IntersectionObserver | null = null
  let isHandlingIntersection = false

  function disconnect() {
    observer?.disconnect()
    observer = null
  }

  async function observe() {
    disconnect()

    if (!options.enabled.value) {
      return
    }

    await nextTick()

    const targetElement = options.target.value
    if (!targetElement || !options.enabled.value) {
      return
    }

    observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (!entry?.isIntersecting || !options.enabled.value || isHandlingIntersection) {
          return
        }

        isHandlingIntersection = true
        Promise.resolve(options.onLoadMore()).finally(() => {
          isHandlingIntersection = false
        })
      },
      {
        root: null,
        rootMargin: options.rootMargin ?? '320px 0px',
        threshold: 0.01,
      },
    )

    observer.observe(targetElement)
  }

  watch([options.target, options.enabled], () => {
    void observe()
  }, { immediate: true })

  onBeforeUnmount(() => {
    disconnect()
  })
}
