# -*- mode: ruby -*-
# vi: set ft=ruby :

# test virtual machines with target operating systems

Vagrant.configure("2") do |config|
  
  config.vm.provider "virtualbox" do |v|
    v.memory = 4096
    v.cpus = 2
  end

  config.vm.provision "file", source: "examples/data-service", destination: "data-service"
  config.vm.provision "file", source: "examples/web-service", destination: "web-service"
  config.vm.provision "file", source: "examples/flask-service", destination: "flask-service"
  config.vm.provision "shell", path: "download.sh"
  config.vm.provision "shell", 
    inline: "sudo mv reqs /usr/bin/"

  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "ubuntu/trusty64"
  end

  config.vm.define "fedora" do |fedora|
    fedora.vm.box = "fedora/28-cloud-base"
  end

  config.vm.define "osx" do |osx|
    osx.vm.box = "AndrewDryga/vagrant-box-osx"
  end

end
