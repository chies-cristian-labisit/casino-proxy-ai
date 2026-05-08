window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/specs/casino-proxy-api.yaml",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout",
    queryConfigEnabled: false,
  })
}
