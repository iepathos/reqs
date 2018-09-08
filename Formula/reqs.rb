class Reqs < Formula
  desc "manage cross-platform system requirements with one tool"
  homepage ""
  url "https://github.com/iepathos/reqs/releases/download/v0.1.1/reqs_0.1.1_Darwin_x86_64.tar.gz", :using => GitHubPrivateRepositoryReleaseDownloadStrategy
  version "0.1.1"
  sha256 "0194e4071b9cae681efd077f0101cd152ff11b4869c919a7f263a1dacbdb8389"

  def install
    bin.install "reqs"
  end
end
