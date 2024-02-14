# Deployment Information

Spawn created applications are able to be deployed in a large number of configurations. The following are some of the most common configurations:

## Single Node

The simplest deployment is a single node. This is the easiest to set up and is the most common for development and testing. It is also the most common for small scale production deployments. This has the following downsides:
- Single point of failure
- No redundancy
- No scalability
- Insecure IBC

TODO: Add some options for easy deployment for single node deployments. If theres some GH action thing to deploy to a cloud that would be ideal.

## Single Node with Horcrux

A single node with Horcrux is a single node deployment with a Horcrux backup. This is a simple way to add redundancy to a single node deployment. This is similar to the  This has the following downsides:
- No redundancy for the chain node
- Still pretty insecure IBC 

TODO: Add a guide for setting this up. this option is relatively hard to do and doesn't really offer a whole ton more but we have it :shrug:

## Full Consensus Network

A full consensus network is a network with a full set of nodes. This is the most secure and scalable deployment. It requires at least 4 nodes to be secure. It is more expensive to run and more complex to set up, but provides much more security for IBC and allows for more complex incetive structures.

TODO: Some of the following: guides for setting up all the nodes yourself, maybe a cloud deployment option. A list of validators who are interested in deploying these types of networks. A talk to us button. 