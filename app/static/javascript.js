'use strict';

const ROVERS = ['spirit', 'curiosity', 'opportunity']

// an object for app state
let State = {}
State.page = 1
State.rover = 0
State.tick = false
State.visible = 0
State.lastY = undefined

document.addEventListener('DOMContentLoaded', function() {
  // cache some dom elements
  State.main = document.querySelector('div.wrapper-main')
  State.nodes = Array.from(document.querySelectorAll('div.wrapper-item'))
  // load in the images
  let photos = Array.from(document.querySelectorAll('img.photo'), (img) => lazyLoad(img))

  // show the first content item
  State.nodes[State.visible].classList.remove('hidden')

  // event handler where the magic happens
  let scrollHandler = (e) => {
    let up
    // for wheel events
    if (e.deltaY && e.deltaY < 0) up = true
    if (e.deltaY && e.deltaY > 0) up = false
    // for touch events
    if (e.type === 'touchmove') {
      if (State.lastY && State.lastY > e.touches[0].clientY) up = true
      if (State.lastY && State.lastY < e.touches[0].clientY) up = false
      State.lastY = e.touches[0].clientY
    }
    if (up && State.visible === 0) return // don't do anything if someone scrolls up right away
    if (!State.tick) {
      window.requestAnimationFrame(function() {
        // manage visible image state
        // first hide the current image
        State.nodes[State.visible].classList.add('hidden')
        // then figure out which one should be shown next
        if (!up) State.visible++
        if (up) State.visible--
        if (State.visible >= State.nodes.length) State.visible = State.nodes.length - 1
        if (State.visible < 0) State.visible = 0
        // now show the correct visible image
        State.nodes[State.visible].classList.remove('hidden')
        // fetch and append new images when scroll is getting close to the end
        if (State.visible >= (State.nodes.length - 10)) {
          let route = makeRoute(ROVERS[State.rover], State.page)
          Util.generator(update, route)
        }
        State.tick = false
      })
    }
    State.tick = true
  }

  ;(function init(){
    State.main.addEventListener('mousewheel', scrollHandler)
    State.main.addEventListener('DOMMouseScroll', scrollHandler)
    State.main.addEventListener('touchmove', scrollHandler)
  })()
})

/**
* Fetch images and metadata and append to the DOM
* @param {string} uri - a generator function to run asynchronously
*/
function* update(uri) {
  try {
    let res = yield fetch(uri) // returns a promise for the response
    let list = yield res.json() // returns a promise for json
    if (!list || list.length < 10) {
      State.rover++
      if (State.rover > ROVERS.length) {
        State.rover = 0
      }
      State.page = 0
      return
    }
    for (let item of list) {
      let node = mkNode(item)
      State.nodes.push(node)
      State.main.append(node)
    }
    State.page++
  } catch (err) {
    console.error(err)
  }
}

/**
* Make new dom nodes with image data
* @param {Object} data - data for one image container div
*
* @return {Object}
*/
function mkNode(data) {
  let div = document.createElement('div')
  div.classList.add('wrapper-item', 'hidden')
  const template = p => Util.template`
  <div class="img-container">
    <img src="${p.img_src}" class="photo"></img>
  </div>
  <div class="metadata">
    <ul>
      <li>image id: ${p.id}</li>
      <li>martian sol: ${p.sol}</li>
      <li>earth date: ${p.earth_date}</li>
      <li>rover: ${p.rover}</li>
      <li>camera: ${p.camera}</li>
    </ul>
  </div>
  `
  div.innerHTML = template(data)
  div.dataset.id = data.id
  return div
}

/**
* Lazy load initial images in rendered markup
* @param {Object} img - an array-like list of <img> DOM nodes in the initial index.html
*
* @return {Object}
*/
function lazyLoad(img) {
  img.src = img.dataset.src
  return img
}

/**
* Convenience function for building the api route
* @param {String} rover - the rover whose data we want
* @param {String} page - the page of data we want
*
* @return {String}
*/
function makeRoute (rover, page) {
  return `/rover/${rover}/page/${page}`
}