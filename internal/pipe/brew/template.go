package brew

import "github.com/goreleaser/goreleaser/pkg/config"

type templateData struct {
	Name                 string
	Desc                 string
	Homepage             string
	Version              string
	License              string
	Caveats              []string
	Plist                string
	PostInstall          []string
	Dependencies         []config.HomebrewDependency
	Conflicts            []string
	Tests                []string
	CustomRequire        string
	CustomBlock          []string
	LinuxPackages        []releasePackage
	MacOSPackages        []releasePackage
	Service              []string
	HasOnlyAmd64MacOsPkg bool
}

type releasePackage struct {
	DownloadURL      string
	SHA256           string
	OS               string
	Arch             string
	DownloadStrategy string
	Install          []string
}

const formulaTemplate = `# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
{{ if .CustomRequire -}}
require_relative "{{ .CustomRequire }}"
{{ end -}}
class {{ .Name }} < Formula
  desc "{{ .Desc }}"
  homepage "{{ .Homepage }}"
  version "{{ .Version }}"
  {{- if .License }}
  license "{{ .License }}"
  {{- end }}
  {{- with .Dependencies }}
  {{ range $index, $element := . }}
  depends_on "{{ .Name }}"
  {{- if .Type }} => :{{ .Type }}{{- else if .Version }} => "{{ .Version }}"{{- end }}
  {{- with .OS }} if OS.{{ . }}?{{- end }}
  {{- end }}
  {{- end -}}

  {{- if and (not .LinuxPackages) .MacOSPackages }}
  depends_on :macos
  {{- end }}
  {{- if and (not .MacOSPackages) .LinuxPackages }}
  depends_on :linux
  {{- end }}
  {{- printf "\n" }}

  {{- if .MacOSPackages }}
  on_macos do
  {{- range $element := .MacOSPackages }}
    {{- if eq $element.Arch "all" }}
    url "{{ $element.DownloadURL }}"
	{{- if .DownloadStrategy }}, using: {{ .DownloadStrategy }}{{- end }}
    sha256 "{{ $element.SHA256 }}"

    def install
      {{- range $index, $element := .Install }}
      {{ . -}}
      {{- end }}
    end
    {{- else if $.HasOnlyAmd64MacOsPkg }}
    url "{{ $element.DownloadURL }}"
	{{- if .DownloadStrategy }}, using: {{ .DownloadStrategy }}{{- end }}
    sha256 "{{ $element.SHA256 }}"

    def install
      {{- range $index, $element := .Install }}
      {{ . -}}
      {{- end }}
    end

    if Hardware::CPU.arm?
      def caveats
        <<~EOS
          The darwin_arm64 architecture is not supported for the {{ $.Name }}
          formula at this time. The darwin_amd64 binary may work in compatibility
          mode, but it might not be fully supported.
        EOS
      end
    end
    {{- else }}
    {{- if eq $element.Arch "amd64" }}
    if Hardware::CPU.intel?
    {{- end }}
    {{- if eq $element.Arch "arm64" }}
    if Hardware::CPU.arm?
    {{- end}}
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, using: {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"

      def install
        {{- range $index, $element := .Install }}
        {{ . -}}
        {{- end }}
      end
    end
    {{- end }}
  {{- end }}
  end
  {{- end }}

  {{- if and .MacOSPackages .LinuxPackages }}{{ printf "\n" }}{{ end }}

  {{- if .LinuxPackages }}
  on_linux do
  {{- range $element := .LinuxPackages }}
    {{- if eq $element.Arch "amd64" }}
    if Hardware::CPU.intel?
    {{- end }}
    {{- if eq $element.Arch "arm" }}
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
    {{- end }}
    {{- if eq $element.Arch "arm64" }}
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    {{- end }}
      url "{{ $element.DownloadURL }}"
	  {{- if .DownloadStrategy }}, using: {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"

      def install
        {{- range $index, $element := .Install }}
        {{ . -}}
        {{- end }}
      end
    end
  {{- end }}
  end
  {{- end }}

  {{- with .Conflicts }}
  {{ range $index, $element := . }}
  conflicts_with "{{ . }}"
  {{- end }}
  {{- end }}

  {{- with .CustomBlock }}
  {{ range $index, $element := . }}
  {{ . }}
  {{- end }}
  {{- end }}

  {{- with .PostInstall }}

  def post_install
    {{- range . }}
    {{ . }}
    {{- end }}
  end
  {{- end -}}

  {{- with .Caveats }}

  def caveats
    <<~EOS
    {{- range $index, $element := . }}
      {{ . -}}
    {{- end }}
    EOS
  end
  {{- end -}}

  {{- with .Plist }}

  plist_options startup: false

  def plist
    <<~EOS
      {{ . }}
    EOS
  end
  {{- end -}}

  {{- with .Service }}

  service do
    {{- range . }}
    {{ . }}
    {{- end }}
  end
  {{- end -}}

  {{- if .Tests }}

  test do
    {{- range $index, $element := .Tests }}
    {{ . -}}
    {{- end }}
  end
  {{- end }}
end
`
