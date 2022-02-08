FROM node:16@sha256:9627b22e9e2e13ee0dbe692d85f3350c082c2b085e4e53114a016143dbc472f2 as builder

WORKDIR /webapp
COPY . .
RUN npm i && npm run build

FROM nginx:1.21.5@sha256:ee89b00528ff4f02f2405e4ee221743ebc3f8e8dd0bfd5c4c20a2fa2aaa7ede3

COPY --from=builder /webapp/dist/ /usr/share/nginx/html