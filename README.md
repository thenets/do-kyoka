[![Docker Repository on Quay](https://quay.io/repository/thenets/do-kyoka/status "Docker Repository on Quay")](https://quay.io/repository/thenets/do-kyoka)

# do-kyoka 💂‍♀️

Auto-update Digital Ocean 🌊 firewall. The easiest way to access all your instances without a VPN.

It gets your public IP address and creates/updates a firewall at Digital Ocean that allows you to join your instances.

## 👨‍💻 Motivation

Digital Ocean is a very good cloud provider for developers who want an easy way to deploy their applications. No fancy features. Just create a droplet (VM), jump into the instance and play with it.

The problem is: many times I don't want to have an instance joining my VPN but I also don't want to have the instance publicly available on the internet.

So, this project is my attempt to solve my problem to allow me to access everything that I want without exposing my application everywhere or deploying my instance in one specific VPC, or joining into my VPN.

> 🔴 This solution IS NOT DESIGNED FOR PRODUCTION. You must use a VPN and a private network to have a secure environment and not expose your application to the internet.

## 📚 Requirements

- Digital Ocean API token with write access (https://cloud.digitalocean.com/account/api/)
- A container runner (docker|podman)

## 🚢 How to use 

The `do-kyoka` creates a firewall at Digital Ocean that allows you to join your instances. Add the corresponding `tag` to your droplet and you will be able to access it.

Container image:
- `quay.io/thenets/do-kyoka:latest`

### 🛠 Environment variables

- `DO_API_TOKEN`: `[required]` Digital Ocean API token with `write` permission. NEVER SHARE IT WITH ANYONE.
- `FIREWALL_NAME`: [default: do-kyoka] The name of the firewall to create/update on Digital Ocean. It's important to notice that ALL CURRENT FIREWALL RULES WILL BE DELETED!
- `FIREWALL_TAG`: [default: do-kyoka] The `tag` to add to your droplet to allow it to join the firewall.

### 💽 Example

```bash
docker run -it \
    --name "do-kyoka-firewall" \
    --restart unless-stopped \
    -e "DO_API_TOKEN=Nxkal9ZWxOMjo6WtU26KLNAnsaW8xQnJjGaT88VkkjyTX8POOdP52Z9XM5K0TM542" \
    -e "FIREWALL_NAME=do-kyoka" \
    -e "FIREWALL_TAG=do-kyoka" \
    quay.io/thenets/do-kyoka:latest
```
