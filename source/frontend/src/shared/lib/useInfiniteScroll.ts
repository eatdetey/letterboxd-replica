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
  let isTargetIntersecting = false

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
        isTargetIntersecting = Boolean(entry?.isIntersecting)

        if (!isTargetIntersecting || !options.enabled.value || isHandlingIntersection) {
          return
        }

        isHandlingIntersection = true
        Promise.resolve(options.onLoadMore()).finally(() => {
          isHandlingIntersection = false

          // Re-arm observer so infinite loading continues even when the sentinel
          // stays inside viewport after the previous batch render.
          if (options.enabled.value && isTargetIntersecting) {
            void observe()
          }
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
