import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client/core'

const httpLinkPages = createHttpLink({
  uri: 'http://localhost:1337/graphql',
})

const cachePages = new InMemoryCache()

export const apolloClientPages = new ApolloClient({
  link: httpLinkPages,
  cache: cachePages,
})

const httpLinkProjects = createHttpLink({
  uri: 'http://localhost:3000/graphql',
})

const cacheProjects = new InMemoryCache()

export const apolloClientProjects = new ApolloClient({
  link: httpLinkProjects,
  cache: cacheProjects,
})