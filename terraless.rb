# This file was generated by GoReleaser. DO NOT EDIT.
class Terraless < Formula
  desc "Terraless is helper to easily deploy Lambda Functions and different projects with Terraform"
  homepage ""
  version "0.3.4"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/Odania-IT/terraless/releases/download/v0.3.4/terraless_darwin_amd64.tar.gz"
    sha256 "c0c2bbccb1cf6e06596162e80748506f6c6b3a0cf5887609475747b4ee49895d"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/Odania-IT/terraless/releases/download/v0.3.4/terraless_linux_amd64.tar.gz"
      sha256 "5ba6961567e8d6a578adc5819edaaf929252b2f50bd4aad37c5f124e71a0e846"
    end
  end

  def install
    bin.install "terraless"
  end
end
