# ADR-002-front-end-framework

Creation Date: 05/01/22

## Context

As part of the development of the Voogle application we have to choose a framework.
We have identified some framework, VueJs and ReactJs

## Problematic

The chosen language must be suited to write Cloud Native application.
The framework must be able to refresh the data instantly and can integrate a video player

## Decision
## Status

Draft

## Options

### Option 1: VueJs

#### Benefits

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
#### Points of attention

### Option 2: ReactJs

#### Benefits
#### Drawbacks
#### Points of attention

## Technical resources
* https://mobiskill.fr/blog/conseils-emploi-tech/vue-js-quels-sont-les-avantages-et-les-inconvenients/