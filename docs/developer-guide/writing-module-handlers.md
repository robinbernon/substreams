---
description: StreamingFast Substreams module handler creation
---

# Module Handler Creation

After the ABI and `Protobuf` Rust code has been generated the handler code needs to be written. The code should be saved into the `src` directory and use the filename `lib.rs.`

{% code title="src/lib.rs" %}
```rust
mod abi;
mod pb;
use hex_literal::hex;
use pb::erc721;
use substreams::{log, store, Hex};
use substreams_ethereum::{pb::eth::v2 as eth, NULL_ADDRESS, Event};

// Bored Ape Yacht Club Contract
const TRACKED_CONTRACT: [u8; 20] = hex!("bc4ca0eda7647a8ab7c2061c2e118a18a936f13d");

/// Extracts transfer events from the contract
#[substreams::handlers::map]
fn block_to_transfers(blk: eth::Block) -> Result<erc721::Transfers, substreams::errors::Error> {
    let mut transfers: Vec<erc721::Transfer> = vec![];
    for trx in blk.transaction_traces {
        transfers.extend(trx.receipt.unwrap().logs.iter().filter_map(|log| {
            if log.address != TRACKED_CONTRACT {
                return None;
            }

            log::debug!("NFT Contract {} invoked", Hex(&TRACKED_CONTRACT));

            if !abi::erc721::events::Transfer::match_log(log) {
                return None;
            }

            let transfer = abi::erc721::events::Transfer::match_and_decode(log).unwrap();

            Some(erc721::Transfer {
                trx_hash: trx.hash.clone(),
                from: transfer.from,
                to: transfer.to,
                token_id: transfer.token_id.low_u64(),
                ordinal: log.block_index as u64,
            })
        }));
    }

    Ok(erc721::Transfers { transfers })
}

// Store the total balance of NFT tokens by address for the specific TRACKED_CONTRACT by holder
#[substreams::handlers::store]
fn nft_state(transfers: erc721::Transfers, s: store::StoreAddInt64) {
    log::info!("NFT state builder");
    for transfer in transfers.transfers {
        if transfer.from != NULL_ADDRESS {
            log::info!("Found a transfer out");

            s.add(transfer.ordinal, generate_key(&transfer.from), -1);
        }

        if transfer.to != NULL_ADDRESS {
            log::info!("Found a transfer in");

            s.add(transfer.ordinal, generate_key(&transfer.to), 1);
        }
    }
}

fn generate_key(holder: &Vec<u8>) -> String {
    return format!("total:{}:{}", Hex(holder), Hex(TRACKED_CONTRACT));
}
```
{% endcode %}

### **Module Handler Breakdown**

Each logical section of the `lib.rs` file is outlined and described in greater detail below.

Import the necessary modules.

```rust
mod abi;
mod pb;
use hex_literal::hex;
use pb::erc721;
use substreams::{log, store, Hex};
use substreams_ethereum::{pb::eth::v2 as eth, NULL_ADDRESS, Event};
```

Store the contract being tracked in the example as a `constant`.

```rust
// Bored Ape Yacht Club Contract
const TRACKED_CONTRACT: [u8; 20] = hex!("bc4ca0eda7647a8ab7c2061c2e118a18a936f13d");
```

Define the `map` module. Here is the module definition from the example Substreams manifest.

```yaml
  - name: block_to_transfers
    kind: map
    initialBlock: 12287507
    inputs:
      - source: sf.ethereum.type.v2.Block
    output:
      type: proto:eth.erc721.v1.Transfers
```

Notice the: `name: block_to_transfers`. This name should correspond to our handler function name.

Also notice, there is one input and one output defined.&#x20;

The input has a type of `sf.ethereum.type.v2.Block`, a standard Ethereum block provided by the `substreams-ethereum` crate.&#x20;

The output is typed as `proto:eth.erc721.v1.Transfers`. This is the custom protobuf definition and is provided by the generated Rust code. Resulting in the following function signature.

```rust
/// Extracts transfers events from the contract
#[substreams::handlers::map]
fn block_to_transfers(blk: eth::Block) -> Result<erc721::Transfers, substreams::errors::Error> {
    ...
}
```

{% hint style="info" %}
_Note:_

_**Rust Macros**_

_Notice the #\[substreams::handlers::map] above the function, this is a rust macro that is provided by the substreams crate. This macro decorates our handler function as a map. There is also a macro used to decorate the handler of kind store:_

_#\[substreams::handlers::store]_
{% endhint %}

The goal of the `map` being built is to extract `ERC721` transfers from a given block.&#x20;

This can be achieved by finding all the `Transfer` events that are emitted by the contract that is currently being tracked. As events are found they will be decoded into Transfer objects.

```rust
/// Extracts transfer events from the contract
#[substreams::handlers::map]
fn block_to_transfers(blk: eth::Block) -> Result<erc721::Transfers, substreams::errors::Error> {
    // variable to store the transfers we find
    let mut transfers: Vec<erc721::Transfer> = vec![];
    // loop through the block's transaction
    for trx in blk.transaction_traces {
        // iterate over the transaction logs
        transfers.extend(trx.receipt.unwrap().logs.iter().filter_map(|log| {
            // verifying that the logs emitted are from the contract we are tracking
            if log.address != TRACKED_CONTRACT {
                return None;
            }

            log::debug!("NFT Contract {} invoked", Hex(&TRACKED_CONTRACT));
            // verify if the log matches a Transfer Event
            if !abi::erc721::events::Transfer::match_log(log) {
                return None;
            }
            
            // decode the event and store it
            let transfer = abi::erc721::events::Transfer::match_and_decode(log).unwrap();
            Some(erc721::Transfer {
                trx_hash: trx.hash.clone(),
                from: transfer.from,
                to: transfer.to,
                token_id: transfer.token_id.low_u64(),
                ordinal: log.block_index as u64,
            })
        }));
    }
    
    // return our list of transfers for the given block
    Ok(erc721::Transfers { transfers })
}
```

Now define the `store` module. As a reminder, here is the module definition from the example Substreams manifest.

```yaml
  - name: nft_state
    kind: store
    initialBlock: 12287507
    updatePolicy: add
    valueType: int64
    inputs:
      - map: block_to_transfers
```

{% hint style="info" %}
_Note: `name: nft_state` will also correspond to the handler function name._
{% endhint %}

The input corresponds to the output of the `block_to_transfers` `map` module typed as `proto:eth.erc721.v1.Transfers`. This is the custom protobuf definition and is provided by the generated Rust code. Resulting in the following function signature.

```rust
/// Store the total balance of NFT tokens for the specific TRACKED_CONTRACT by holder
#[substreams::handlers::store]
fn nft_state(transfers: erc721::Transfers, s: store::StoreAddInt64) {
    ...
}
```

{% hint style="info" %}
_Note: the `store` will always receive itself as its own last input._&#x20;
{% endhint %}

In this example the `store` module uses an `updatePolicy` set to `add` and a `valueType set` to `int64` yielding a writable store typed as `StoreAddInt64`.

{% hint style="info" %}
_Note: ****_&#x20;

_**Store Types**_

_The last parameter of a `store` module function should always be the writable store itself._&#x20;

_The type of the writable store is based on the `store` module `updatePolicy` and `valueType`_.&#x20;
{% endhint %}

The goal of the `store` in this example is to track a holder's current NFT count for the contract supplied. This tracking is achieved through the analyzation of transfers.

**Transfer in Detail**

If the transfer's `from` address field contains the null address (`0x0000000000000000000000000000000000000000`), and the `to` address field is not the null address, the `to` address field is minting a token, so the count should be incremented.

If the transfer's `from` address field is not the null address, _and_ the `to` address field is the null address, the `from` address field is burning a token, so the count should be decremented.

If the `from` address field and the `to` address field is not a null address, the count should be decremented of the `from` address, and increment the count of the `to` address for basic transfers.

### Store Concepts

When writing to a store, there are three concepts to consider that include `ordinal`, `key` and `value`. Additional information for each is provided below.

#### Ordinal

Ordinal represents the order in which the `store` operations will be applied.&#x20;

The `store` handler will be called once per `block.`&#x20;

During execution, the `add` operation may be called multiple times, for multiple reasons, such as finding a relevant event or seeing a call that triggered a method call.&#x20;

Blockchain execution models are linear. Operations to add must be added linearly and deterministically.

When an ordinal is specified the order of execution is guaranteed. For one execution of the `store` handler for given inputs, in this example a list of transfers, the code will emit the same number of `add` calls and ordinal values.

#### Key

Stores are [key/value stores](https://en.wikipedia.org/wiki/Key%E2%80%93value\_database). Care needs to be taken when crafting a key to ensure that it is unique _and flexible_.&#x20;

In the example, if the `generate_key` function would simply return a key that is the `TRACKED_CONTRACT` address it would not be unique between different token holders.&#x20;

If the `generate_key` function returned a key containing only the holder's address it would be unique amongst holders. Issues would be encountered however when attempting to track multiple contracts.

#### Value

The value being stored. The type is dependent on the store type being used.

```rust
/// Store the total balance of NFT tokens for the specific TRACKED_CONTRACT by holder
#[substreams::handlers::store]
fn nft_state(transfers: erc721::Transfers, s: store::StoreAddInt64) {
    log::info!("NFT state builder");
    // iterate over the transfers event
    for transfer in transfers.transfers {
        // check if the from address field is not the NULL address
        if transfer.from != NULL_ADDRESS {
            log::info!("Found a transfer out");
            // decrement the count
            s.add(transfer.ordinal, generate_key(&transfer.from), -1);
        }
        // check if the to address field is not the NULL address
        if transfer.to != NULL_ADDRESS {
            log::info!("Found a transfer in");
            // increment the count
            s.add(transfer.ordinal, generate_key(&transfer.to), 1);
        }
    }
}

fn generate_key(holder: &Vec<u8>) -> String {
    return format!("total:{}:{}", Hex(holder), Hex(TRACKED_CONTRACT));
}
```

### Summary

Both handler functions have been written.&#x20;

One handler function for extracting transfers that are of interest, and a second to store the token count per recipient.&#x20;

Build Substreams to continue the setup process.&#x20;

```
cargo build --target wasm32-unknown-unknown --release
```

The next step is to run Substreams with all of the changes made using the code that's been generated.
