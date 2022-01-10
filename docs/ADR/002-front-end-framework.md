# ADR-002-front-end-framework

* Creation Date: 05/01/22

## Context

As part of the development of the Voogle application we have to choose a frontend framework.
We have identified some framework, VueJs and ReactJs

## Problematic

The chosen language must be suited to write Cloud Native application.
The framework must be able to refresh the data instantly and can integrate a video player

## Decision

We chose to use VueJs. We need a simple framework to develop quickly. Moreover our application is not complex.

## Options

### Option 1: VueJs

*"Vue.js is an open-source model–view–viewmodel front end JavaScript framework for building user interfaces and single-page applications. It was created by Evan You, and is maintained by him and the rest of the active core team members"*
[Wikipedia](https://en.wikipedia.org/wiki/Vue.js)

#### Benefits

* Opensource
* Simple approach: javascript and html
* Data binding
* Bidirectional communication architecture: MVVM and Virtual DOM
* Detailed and complete documentation
* Flexibilité: run with browser
* Simple integration: integration with existing. No need to start form scratch
* Small size: framework very small 18-21 Ko and faster than Angular.JS and React.JS and Ember.JS.
* Versatility in terms of application size: Vue companion libraries are officially supported and are kept up-to-date with the main library while Redux is an unofficial react extension (redux -> vuex)

#### Drawbacks

* Lack of scalability
* Lack pluging compared to Angular et React
* Little community

### Option 2: ReactJs

*"React (also known as React.js or ReactJS) is a free and open-source front-end JavaScript library for building user interfaces based on UI components. It is maintained by Meta (formerly Facebook) and a community of individual developers and companies. React can be used as a base in the development of single-page or mobile applications. However, React is only concerned with state management and rendering that state to the DOM, so creating React applications usually requires the use of additional libraries for routing, as well as certain client-side functionality."*
[Wikipedia](https://en.wikipedia.org/wiki/React_(JavaScript_library))

#### Benefits

* Opensource
* Simplicity: component approach
* JSX ?
* Reusability
* Optimized for search engine optimization -> SEO of google
* Large community
* Many plugin
* Flexibility et responsiveness
* Reduction of the number of operations on the DOM

#### Drawbacks
* JSX ?
* Very little official documentation
* React is unopinionated
* Long to master, ReactJS requires a deep knowledge of how to integrate the user interface into the MVC framework.

### Options diff

If you are looking for something with a faster development time and an application with less complexity, Vue is a good choice, but React has more diversity, as it allows the development of very complex applications.

## Technical resources
* https://mobiskill.fr/blog/conseils-emploi-tech/vue-js-quels-sont-les-avantages-et-les-inconvenients/