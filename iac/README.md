# Notes

- This is a basic Linux VM, which I have created using terraform. 
- It runs ubuntu 22 in the us-west-2 region on a t3.micro instance.

## Steps before use

- Install Terrform
- Setup a basic AWS account and configure an active session in your terminal
- Add a key pair to your aws account (this is free). AWS can help you do this.
- In the hg1 file, I have used my key, "gen-aws". You have to change it to point to your key, so that you can connect to the VM created from your computer. Make sure you have the private key.

In the `iac` directory(or wherever you have the hg1.tf file) run,

```
terraform init 

terraform plan

terraform apply
```

You will get the Public IP, Private IP and Key Name of the instance.

- You can now SSH into the VM, using this command.

```
ssh -i <path-to-private-key> ubuntu@<public-ip>
```

- In order for the challenges to work, you need another port open, so that the protohackers clients can connect with your VM. You have to edit the default security group, to allow traffic to this port. (I'll let you figure this one out as well)

Happy Hacking :)