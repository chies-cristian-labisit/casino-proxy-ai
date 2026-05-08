window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/specs/admin-api-spec.yaml",
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
